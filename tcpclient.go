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

package mr2

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"time"

	"github.com/gogo/protobuf/proto"
	cache "github.com/patrickmn/go-cache"
)

//TCPClient .
type TCPClient struct {
	Client  *Client
	TCPConn *net.TCPConn
	Cache   *cache.Cache
	Done    chan byte
	Data    chan []byte
	Error   chan error
}

// NewTCPClient .
func NewTCPClient(c *Client) (*TCPClient, error) {
	tmp, err := Dial.Dial("tcp", c.Server)
	if err != nil {
		return nil, err
	}
	conn := tmp.(*net.TCPConn)
	if err := conn.SetKeepAlivePeriod(time.Duration(60) * time.Second); err != nil {
		conn.Close()
		return nil, err
	}
	if err := conn.SetDeadline(time.Now().Add(time.Duration(10) * time.Second)); err != nil {
		conn.Close()
		return nil, err
	}
	tmp1, err := c.Ckv.Encrypt("Mr.2", "TCPHello")
	if err != nil {
		conn.Close()
		return nil, err
	}
	p := &TCPHello{
		Domain:      c.ServerDomain,
		Port:        c.ServerPort,
		TCPTimeout:  c.TCPTimeout,
		TCPDeadline: c.TCPDeadline,
		Key:         tmp1,
	}
	b, err := proto.Marshal(p)
	if err != nil {
		conn.Close()
		return nil, err
	}
	bb := make([]byte, 2)
	binary.BigEndian.PutUint16(bb, uint16(len(b)))
	if _, err := conn.Write(append(bb, b...)); err != nil {
		conn.Close()
		return nil, err
	}
	tc := &TCPClient{
		Client:  c,
		Cache:   cache.New(cache.NoExpiration, cache.NoExpiration),
		TCPConn: conn,
		Done:    make(chan byte),
		Data:    make(chan []byte),
		Error:   make(chan error),
	}
	return tc, nil
}

// Run .
func (c *TCPClient) Run() error {
	defer close(c.Done)
	defer c.TCPConn.Close()

	go c.Ping()
	go c.Read()

	for {
		select {
		case b := <-c.Data:
			if err := c.TCPConn.SetDeadline(time.Now().Add(time.Duration(10) * time.Second)); err != nil {
				return err
			}
			if _, err := c.TCPConn.Write(b); err != nil {
				return err
			}
		case err := <-c.Error:
			return err
		}
	}
	return nil
}

// Stop .
func (c *TCPClient) Stop() {
	select {
	case <-c.Done:
		return
	case c.Error <- nil:
	}
}

// Ping server .
func (c *TCPClient) Ping() {
	for {
		time.Sleep(5 * time.Second)
		p := &PingPong{}
		b, err := proto.Marshal(p)
		if err != nil {
			select {
			case <-c.Done:
				return
			case c.Error <- err:
			}
			return
		}
		bb := make([]byte, 2)
		binary.BigEndian.PutUint16(bb, uint16(len(b)))
		select {
		case <-c.Done:
			return
		case c.Data <- append(append([]byte{0x00}, bb...), b...):
		}
	}
}

// Read data from server.
func (c *TCPClient) Read() {
	for {
		if err := c.TCPConn.SetDeadline(time.Now().Add(time.Duration(10) * time.Second)); err != nil {
			select {
			case <-c.Done:
				return
			case c.Error <- err:
			}
			return
		}
		b := make([]byte, 3)
		if _, err := io.ReadFull(c.TCPConn, b); err != nil {
			select {
			case <-c.Error:
				return
			case c.Error <- err:
			}
			return
		}
		k := b[0]
		i := int(binary.BigEndian.Uint16(b[1:]))
		b = make([]byte, i)
		if _, err := io.ReadFull(c.TCPConn, b); err != nil {
			select {
			case <-c.Done:
				return
			case c.Error <- err:
			}
			return
		}
		if k == 0x00 {
			continue
		}
		if k == 0x02 {
			p := &TCPPacket{}
			if err := proto.Unmarshal(b, p); err != nil {
				select {
				case <-c.Done:
					return
				case c.Error <- err:
				}
				return
			}
			i, ok := c.Cache.Get(p.Address)
			if ok {
				c1 := i.(*net.TCPConn)
				c.Cache.Delete(p.Address)
				c1.Close()
				return
			}
			continue
		}
		if k == 0x01 {
			p := &TCPPacket{}
			if err := proto.Unmarshal(b, p); err != nil {
				select {
				case <-c.Done:
					return
				case c.Error <- err:
				}
				return
			}
			go func(p *TCPPacket) {
				i, ok := c.Cache.Get(p.Address)
				if ok {
					c1 := i.(*net.TCPConn)
					if _, err := c1.Write(p.Data); err != nil {
						return
					}
					return
				}
				tmp, err := Dial.Dial("tcp", c.Client.ClientServer)
				if err != nil {
					log.Println(err)
					return
				}
				c1 := tmp.(*net.TCPConn)
				defer c1.Close()
				c.Cache.Set(p.Address, c1, cache.DefaultExpiration)
				defer c.Cache.Delete(p.Address)
				if c.Client.TCPTimeout != 0 {
					if err := c1.SetKeepAlivePeriod(time.Duration(c.Client.TCPTimeout) * time.Second); err != nil {
						log.Println(err)
						return
					}
				}
				if c.Client.TCPDeadline != 0 {
					if err := c1.SetDeadline(time.Now().Add(time.Duration(c.Client.TCPDeadline) * time.Second)); err != nil {
						log.Println(err)
						return
					}
				}
				if _, err := c1.Write(p.Data); err != nil {
					log.Println(err)
					return
				}
				var bf [1024 * 2]byte
				for {
					if c.Client.TCPDeadline != 0 {
						if err := c1.SetDeadline(time.Now().Add(time.Duration(c.Client.TCPDeadline) * time.Second)); err != nil {
							return
						}
					}
					i, err := c1.Read(bf[:])
					if err != nil {
						return
					}
					p := &TCPPacket{
						Address: p.Address,
						Data:    bf[0:i],
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
					bb := make([]byte, 2)
					binary.BigEndian.PutUint16(bb, uint16(len(b)))
					select {
					case <-c.Done:
						return
					case c.Data <- append(append([]byte{0x01}, bb...), b...):
					}
				}
			}(p)
		}
	}
}
