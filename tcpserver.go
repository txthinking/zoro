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
	"errors"
	"io"
	"net"
	"strconv"
	"time"

	"github.com/gogo/protobuf/proto"
	cache "github.com/patrickmn/go-cache"
)

// TCPServer .
type TCPServer struct {
	Server      *Server
	Cache       *cache.Cache
	TCPConn     *net.TCPConn
	TCPListen   *net.TCPListener
	TCPTimeout  int64
	TCPDeadline int64
	Done        chan byte
	Data        chan []byte
	Error       chan error
}

// NewTCPServer. Some cases return nil, nil
func NewTCPServer(s *Server, c *net.TCPConn) (*TCPServer, error) {
	if err := c.SetKeepAlivePeriod(time.Duration(60) * time.Second); err != nil {
		return nil, nil
	}
	if err := c.SetDeadline(time.Now().Add(time.Duration(10) * time.Second)); err != nil {
		return nil, nil
	}
	b := make([]byte, 2)
	if _, err := io.ReadFull(c, b); err != nil {
		return nil, nil
	}
	i := int(binary.BigEndian.Uint16(b))
	b = make([]byte, i)
	if _, err := io.ReadFull(c, b); err != nil {
		return nil, nil
	}
	h := &TCPHello{}
	if err := proto.Unmarshal(b, h); err != nil {
		return nil, nil
	}
	if h.Port == 0 {
		return nil, errors.New(c.RemoteAddr().String() + " missed port")
	}
	if len(s.PortCkv) == 0 {
		tmp, err := s.Ckv.Decrypt(h.Key, "Mr.2", 3*60)
		if err != nil || tmp != "TCPHello" {
			return nil, errors.New(c.RemoteAddr().String() + " Hacking")
		}
	}
	if len(s.PortCkv) != 0 {
		ckv, ok := s.PortCkv[h.Port]
		if !ok {
			return nil, errors.New(c.RemoteAddr().String() + " try to open not allowed TCP port: " + strconv.FormatInt(h.Port, 10))
		}
		tmp, err := ckv.Decrypt(h.Key, "Mr.2", 3*60)
		if err != nil || tmp != "TCPHello" {
			return nil, errors.New(c.RemoteAddr().String() + " Hacking")
		}
	}
	taddr, err := net.ResolveTCPAddr("tcp", ":"+strconv.FormatInt(h.Port, 10))
	if err != nil {
		return nil, err
	}
	l, err := net.ListenTCP("tcp", taddr)
	if err != nil {
		return nil, err
	}
	return &TCPServer{
		Server:      s,
		Cache:       cache.New(cache.NoExpiration, cache.NoExpiration),
		TCPConn:     c,
		TCPListen:   l,
		TCPTimeout:  h.TCPTimeout,
		TCPDeadline: h.TCPDeadline,
		Done:        make(chan byte),
		Data:        make(chan []byte),
		Error:       make(chan error),
	}, nil
}

// ListenAndServe .
func (s *TCPServer) ListenAndServe() error {
	defer close(s.Done)
	defer s.TCPListen.Close()
	go s.Accept()
	go s.Read()
	for {
		select {
		case b := <-s.Data:
			if err := s.TCPConn.SetDeadline(time.Now().Add(time.Duration(10) * time.Second)); err != nil {
				return err
			}
			if _, err := s.TCPConn.Write(b); err != nil {
				return err
			}
		case err := <-s.Error:
			return err
		}
	}
	return nil
}

// Shutdown .
func (s *TCPServer) Shutdown() {
	select {
	case <-s.Done:
		return
	case s.Error <- nil:
	}
}

// Accept consumer.
func (s *TCPServer) Accept() {
	for {
		c1, err := s.TCPListen.AcceptTCP()
		if err != nil {
			select {
			case <-s.Done:
				return
			case s.Error <- err:
			}
			return
		}
		if s.TCPTimeout != 0 {
			if err := c1.SetKeepAlivePeriod(time.Duration(s.TCPTimeout) * time.Second); err != nil {
				c1.Close()
				continue
			}
		}
		if s.TCPDeadline != 0 {
			if err := c1.SetDeadline(time.Now().Add(time.Duration(s.TCPDeadline) * time.Second)); err != nil {
				c1.Close()
				continue
			}
		}
		s.Cache.Set(c1.RemoteAddr().String(), c1, cache.DefaultExpiration)
		go func(c1 *net.TCPConn) {
			defer c1.Close()
			defer s.Cache.Delete(c1.RemoteAddr().String())
			defer func() {
				p := &TCPPacket{
					Address: c1.RemoteAddr().String(),
				}
				b, err := proto.Marshal(p)
				if err != nil {
					select {
					case <-s.Done:
						return
					case s.Error <- err:
					}
					return
				}
				bb := make([]byte, 2)
				binary.BigEndian.PutUint16(bb, uint16(len(b)))
				select {
				case <-s.Done:
					return
				case s.Data <- append(append([]byte{0x02}, bb...), b...):
				}
			}()
			var bf [1024 * 2]byte
			for {
				if s.TCPDeadline != 0 {
					if err := c1.SetDeadline(time.Now().Add(time.Duration(s.TCPDeadline) * time.Second)); err != nil {
						return
					}
				}
				i, err := c1.Read(bf[:])
				if err != nil {
					return
				}
				p := &TCPPacket{
					Address: c1.RemoteAddr().String(),
					Data:    bf[0:i],
				}
				b, err := proto.Marshal(p)
				if err != nil {
					select {
					case <-s.Done:
						return
					case s.Error <- err:
					}
					return
				}
				bb := make([]byte, 2)
				binary.BigEndian.PutUint16(bb, uint16(len(b)))
				select {
				case <-s.Done:
					return
				case s.Data <- append(append([]byte{0x01}, bb...), b...):
				}
			}
		}(c1)
	}
}

// Read data from client.
func (s *TCPServer) Read() {
	for {
		if err := s.TCPConn.SetDeadline(time.Now().Add(time.Duration(10) * time.Second)); err != nil {
			select {
			case <-s.Done:
				return
			case s.Error <- nil:
			}
			return
		}
		b := make([]byte, 3)
		if _, err := io.ReadFull(s.TCPConn, b); err != nil {
			select {
			case <-s.Done:
				return
			case s.Error <- nil:
			}
			return
		}
		k := b[0]
		i := int(binary.BigEndian.Uint16(b[1:]))
		b = make([]byte, i)
		if _, err := io.ReadFull(s.TCPConn, b); err != nil {
			select {
			case <-s.Done:
				return
			case s.Error <- nil:
			}
			return
		}
		if k == 0x00 {
			p := &PingPong{}
			if err := proto.Unmarshal(b, p); err != nil {
				select {
				case <-s.Done:
					return
				case s.Error <- err:
				}
				return
			}
			b, err := proto.Marshal(&PingPong{})
			if err != nil {
				select {
				case <-s.Done:
					return
				case s.Error <- err:
				}
				return
			}
			bb := make([]byte, 2)
			binary.BigEndian.PutUint16(bb, uint16(len(b)))
			select {
			case <-s.Done:
				return
			case s.Data <- append(append([]byte{0x00}, bb...), b...):
			}
			continue
		}
		if k == 0x01 {
			p := &TCPPacket{}
			if err := proto.Unmarshal(b, p); err != nil {
				select {
				case <-s.Done:
					return
				case s.Error <- err:
				}
				return
			}
			i, ok := s.Cache.Get(p.Address)
			if !ok {
				continue
			}
			c1 := i.(*net.TCPConn)
			if _, err := c1.Write(p.Data); err != nil {
				continue
			}
			continue
		}
	}
}
