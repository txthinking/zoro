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

package https

import (
	"crypto/tls"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"
	cache "github.com/patrickmn/go-cache"
	"github.com/txthinking/mr2"
)

type DomainData struct {
	Domain string
	Data   []byte
}

// TLSServer .
type TLSServer struct {
	HTTPSServer *HTTPSServer
	Cache       *cache.Cache
	TCPListen   *net.TCPListener
	Done        chan byte
	DomainData  chan DomainData
	Error       chan error
}

// NewTLSServer .
func NewTLSServer(s *HTTPSServer) (*TLSServer, error) {
	l, err := net.ListenTCP("tcp", s.TLSAddr)
	if err != nil {
		return nil, err
	}
	return &TLSServer{
		HTTPSServer: s,
		Cache:       cache.New(cache.NoExpiration, cache.NoExpiration),
		TCPListen:   l,
		Done:        make(chan byte),
		DomainData:  make(chan DomainData),
		Error:       make(chan error),
	}, nil
}

// ListenAndServe .
func (s *TLSServer) ListenAndServe() error {
	defer close(s.Done)
	defer s.TCPListen.Close()
	go s.Accept()
	for {
		select {
		case d := <-s.DomainData:
			i, ok := s.Cache.Get("domain:" + d.Domain)
			if !ok {
				continue
			}
			c := i.(*net.TCPConn)
			if err := c.SetDeadline(time.Now().Add(time.Duration(10) * time.Second)); err != nil {
				continue
			}
			if _, err := c.Write(d.Data); err != nil {
				continue
			}
		case err := <-s.Error:
			return err
		}
	}
	return nil
}

// Shutdown .
func (s *TLSServer) Shutdown() {
	select {
	case <-s.Done:
		return
	case s.Error <- nil:
	}
}

// Accept consumer.
func (s *TLSServer) Accept() {
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
		if s.HTTPSServer.TLSTimeout != 0 {
			if err := c1.SetKeepAlivePeriod(time.Duration(s.HTTPSServer.TLSTimeout) * time.Second); err != nil {
				c1.Close()
				continue
			}
		}
		if s.HTTPSServer.TLSDeadline != 0 {
			if err := c1.SetDeadline(time.Now().Add(time.Duration(s.HTTPSServer.TLSDeadline) * time.Second)); err != nil {
				c1.Close()
				continue
			}
		}
		tc := tls.Server(c1, s.HTTPSServer.TLSConfig)
		if err := tc.Handshake(); err != nil {
			c1.Close()
			continue
		}
		cs := tc.ConnectionState()
		if cs.ServerName == "" {
			log.Println(tc.RemoteAddr().String() + " no SNI")
			tc.Close()
			continue
		}
		s.Cache.Set(tc.RemoteAddr().String(), tc, cache.DefaultExpiration)
		go func(tc *tls.Conn) {
			defer func() {
				p := &mr2.TCPPacket{
					Address: tc.RemoteAddr().String(),
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
				d := DomainData{
					Domain: strings.TrimSuffix(cs.ServerName, "."+s.HTTPSServer.Domain),
					Data:   append(append([]byte{0x02}, bb...), b...),
				}
				select {
				case <-s.Done:
					return
				case s.DomainData <- d:
				}
				s.Cache.Delete(tc.RemoteAddr().String())
				tc.Close()
			}()
			var bf [1024 * 2]byte
			for {
				if s.HTTPSServer.TLSDeadline != 0 {
					if err := tc.SetDeadline(time.Now().Add(time.Duration(s.HTTPSServer.TLSDeadline) * time.Second)); err != nil {
						return
					}
				}
				i, err := tc.Read(bf[:])
				if err != nil {
					return
				}
				p := &mr2.TCPPacket{
					Address: tc.RemoteAddr().String(),
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
				d := DomainData{
					Domain: strings.TrimSuffix(cs.ServerName, "."+s.HTTPSServer.Domain),
					Data:   append(append([]byte{0x01}, bb...), b...),
				}
				select {
				case <-s.Done:
					return
				case s.DomainData <- d:
				}
			}
		}(tc)
	}
}

// HandleClient .
func (s *TLSServer) HandleClient(c *net.TCPConn) error {
	if err := c.SetKeepAlivePeriod(time.Duration(60) * time.Second); err != nil {
		return nil
	}
	if err := c.SetDeadline(time.Now().Add(time.Duration(10) * time.Second)); err != nil {
		return nil
	}
	b := make([]byte, 2)
	if _, err := io.ReadFull(c, b); err != nil {
		return nil
	}
	i := int(binary.BigEndian.Uint16(b))
	b = make([]byte, i)
	if _, err := io.ReadFull(c, b); err != nil {
		return nil
	}
	h := &mr2.TCPHello{}
	if err := proto.Unmarshal(b, h); err != nil {
		return nil
	}
	if h.Domain == "" {
		return errors.New(c.RemoteAddr().String() + " missed domain")
	}
	if len(s.HTTPSServer.DomainCkv) == 0 {
		tmp, err := s.HTTPSServer.Ckv.Decrypt(h.Key, "Mr.2", 3*60)
		if err != nil || tmp != "TCPHello" {
			return errors.New(c.RemoteAddr().String() + " Hacking")
		}
	}
	if len(s.HTTPSServer.DomainCkv) != 0 {
		ckv, ok := s.HTTPSServer.DomainCkv[h.Domain]
		if !ok {
			return errors.New(c.RemoteAddr().String() + " try to open not allowed domain: " + h.Domain)
		}
		tmp, err := ckv.Decrypt(h.Key, "Mr.2", 3*60)
		if err != nil || tmp != "TCPHello" {
			return errors.New(c.RemoteAddr().String() + " Hacking")
		}
	}
	_, ok := s.Cache.Get("domain:" + h.Domain)
	if ok {
		return errors.New(c.RemoteAddr().String() + " try to open domain: " + h.Domain + ", but domain is being used")
	}
	s.Cache.Set("domain:"+h.Domain, c, cache.DefaultExpiration)
	defer s.Cache.Delete("domain:" + h.Domain)
	for {
		if err := c.SetDeadline(time.Now().Add(time.Duration(10) * time.Second)); err != nil {
			return nil
		}
		b := make([]byte, 3)
		if _, err := io.ReadFull(c, b); err != nil {
			return nil
		}
		k := b[0]
		i := int(binary.BigEndian.Uint16(b[1:]))
		b = make([]byte, i)
		if _, err := io.ReadFull(c, b); err != nil {
			return nil
		}
		if k == 0x00 {
			p := &mr2.PingPong{}
			if err := proto.Unmarshal(b, p); err != nil {
				return err
			}
			b, err := proto.Marshal(&mr2.PingPong{})
			if err != nil {
				return err
			}
			bb := make([]byte, 2)
			binary.BigEndian.PutUint16(bb, uint16(len(b)))
			d := DomainData{
				Domain: h.Domain,
				Data:   append(append([]byte{0x00}, bb...), b...),
			}
			select {
			case <-s.Done:
				return nil
			case s.DomainData <- d:
			}
			continue
		}
		if k == 0x01 {
			p := &mr2.TCPPacket{}
			if err := proto.Unmarshal(b, p); err != nil {
				return err
			}
			i, ok := s.Cache.Get(p.Address)
			if !ok {
				continue
			}
			c1 := i.(*tls.Conn)
			if _, err := c1.Write(p.Data); err != nil {
				continue
			}
			continue
		}
		if k == 0x02 {
			p := &mr2.TCPPacket{}
			if err := proto.Unmarshal(b, p); err != nil {
				return err
			}
			i, ok := s.Cache.Get(p.Address)
			if ok {
				c1 := i.(*tls.Conn)
				s.Cache.Delete(p.Address)
				c1.Close()
			}
			continue
		}
	}
	return nil
}
