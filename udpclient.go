// Copyright (c) 2019-present Cloud <cloud@txthinking.com>
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of version 3 of the GNU General Public
// License as published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
// General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package zoro

import (
	"errors"
	"log"
	"net"
	"time"

	"github.com/gogo/protobuf/proto"
	cache "github.com/patrickmn/go-cache"
)

// UDPClient .
type UDPClient struct {
	Client  *Client
	UDPConn *net.UDPConn
	Cache   *cache.Cache
	Done    chan byte
	Data    chan []byte
	Error   chan error
}

// NewUDPClient .
func NewUDPClient(c *Client) (*UDPClient, error) {
	tmp, err := Dial.Dial("udp", c.Server)
	if err != nil {
		return nil, err
	}
	conn := tmp.(*net.UDPConn)
	tmp1, err := c.Ckv.Encrypt("Mr.2", "UDPPacket")
	if err != nil {
		conn.Close()
		return nil, err
	}
	p := &UDPPacket{
		Port: c.ServerPort,
		Key:  tmp1,
	}
	bb, err := proto.Marshal(p)
	if err != nil {
		conn.Close()
		return nil, err
	}
	if _, err := conn.Write(bb); err != nil {
		conn.Close()
		return nil, err
	}
	var b [65536]byte
	i, err := conn.Read(b[:])
	if err != nil {
		conn.Close()
		return nil, err
	}
	p = &UDPPacket{}
	if err := proto.Unmarshal(b[0:i], p); err != nil {
		conn.Close()
		return nil, err
	}
	if p.Address != "0" {
		conn.Close()
		return nil, errors.New(p.Address)
	}
	uc := &UDPClient{
		Client:  c,
		UDPConn: conn,
		Cache:   cache.New(cache.NoExpiration, cache.NoExpiration),
		Done:    make(chan byte),
		Data:    make(chan []byte),
		Error:   make(chan error),
	}
	return uc, nil
}

// Run .
func (c *UDPClient) Run() error {
	defer close(c.Done)
	defer c.UDPConn.Close()
	go c.Ping()
	go c.Read()
	for {
		select {
		case b := <-c.Data:
			if _, err := c.UDPConn.Write(b); err != nil {
				return err
			}
		case err := <-c.Error:
			return err
		}
	}
	return nil
}

// Stop .
func (c *UDPClient) Stop() {
	select {
	case <-c.Done:
		return
	case c.Error <- nil:
	}
}

// Ping server .
func (c *UDPClient) Ping() {
	for {
		time.Sleep(5 * time.Second)
		p := &UDPPacket{
			Port:    c.Client.ServerPort,
			Address: "0",
		}
		bb, err := proto.Marshal(p)
		if err != nil {
			select {
			case <-c.Done:
				return
			case c.Error <- err:
			}
			return
		}
		select {
		case <-c.Done:
			return
		case c.Data <- bb:
		}
	}
}

// Read data from server.
func (c *UDPClient) Read() {
	var b [65536]byte
	for {
		i, err := c.UDPConn.Read(b[:])
		if err != nil {
			select {
			case <-c.Done:
				return
			case c.Error <- err:
			}
			return
		}
		p := &UDPPacket{}
		if err := proto.Unmarshal(b[0:i], p); err != nil {
			continue
		}
		it, ok := c.Cache.Get(p.Address)
		if ok {
			c1 := it.(*net.UDPConn)
			if c.Client.UDPDeadline != 0 {
				if err := c1.SetDeadline(time.Now().Add(time.Duration(c.Client.UDPDeadline) * time.Second)); err != nil {
					continue
				}
			}
			if _, err := c1.Write(p.Data); err != nil {
				continue
			}
			continue
		}
		tmp, err := Dial.Dial("udp", c.Client.ClientServer)
		if err != nil {
			log.Println(err)
			continue
		}
		c1 := tmp.(*net.UDPConn)
		if c.Client.UDPDeadline != 0 {
			if err := c1.SetDeadline(time.Now().Add(time.Duration(c.Client.UDPDeadline) * time.Second)); err != nil {
				log.Println(err)
				c1.Close()
				continue
			}
		}
		if _, err := c1.Write(p.Data); err != nil {
			log.Println(err)
			c1.Close()
			return
		}
		c.Cache.Set(p.Address, c1, cache.DefaultExpiration)
		go func(p *UDPPacket, c1 *net.UDPConn) {
			defer c1.Close()
			defer c.Cache.Delete(p.Address)
			var b [65536]byte
			for {
				if c.Client.UDPDeadline != 0 {
					if err := c1.SetDeadline(time.Now().Add(time.Duration(c.Client.UDPDeadline) * time.Second)); err != nil {
						return
					}
				}
				i, err := c1.Read(b[:])
				if err != nil {
					return
				}
				p := &UDPPacket{
					Port:    c.Client.ServerPort,
					Address: p.Address,
					Data:    b[0:i],
				}
				b, err := proto.Marshal(p)
				if err != nil {
					select {
					case <-c.Done:
						return
					case c.Error <- err:
					}
					return
				}
				select {
				case <-c.Done:
					return
				case c.Data <- b:
				}
			}
		}(p, c1)
	}
}
