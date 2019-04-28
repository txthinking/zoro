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

package main

import (
	"net"
	"strconv"
	"time"

	"github.com/gogo/protobuf/proto"
)

// UDPServer
type UDPServer struct {
	Server  *Server
	UDPConn *net.UDPConn
	Addr    *net.UDPAddr
}

// NewUDPServer
func NewUDPServer(s *Server, port int64, addr *net.UDPAddr) (*UDPServer, error) {
	uaddr, err := net.ResolveUDPAddr("udp", ":"+strconv.FormatInt(port, 10))
	if err != nil {
		return nil, err
	}
	c1, err := net.ListenUDP("udp", uaddr)
	if err != nil {
		p := &UDPPacket{
			Address: err.Error(),
		}
		b, err1 := proto.Marshal(p)
		if err1 != nil {
			return nil, err1
		}
		if _, err := s.UDPConn.WriteToUDP(b, addr); err != nil {
			return nil, err
		}
		return nil, err
	}
	if err := c1.SetDeadline(time.Now().Add(time.Duration(10) * time.Second)); err != nil {
		c1.Close()
		p := &UDPPacket{
			Address: err.Error(),
		}
		b, err := proto.Marshal(p)
		if err != nil {
			return nil, err
		}
		if _, err := s.UDPConn.WriteToUDP(b, addr); err != nil {
			return nil, err
		}
		return nil, err
	}
	p := &UDPPacket{
		Address: "0",
	}
	b, err := proto.Marshal(p)
	if err != nil {
		c1.Close()
		return nil, err
	}
	if _, err := s.UDPConn.WriteToUDP(b, addr); err != nil {
		c1.Close()
		return nil, err
	}
	return &UDPServer{
		Server:  s,
		UDPConn: c1,
		Addr:    addr,
	}, nil
}

// ListenAndServe
func (s *UDPServer) ListenAndServe() error {
	defer s.UDPConn.Close()
	for {
		b := make([]byte, 65536)
		i, a, err := s.UDPConn.ReadFromUDP(b)
		if err != nil {
			return nil
		}
		p := &UDPPacket{
			Address: a.String(),
			Data:    b[0:i],
		}
		b, err = proto.Marshal(p)
		if err != nil {
			return err
		}
		if _, err := s.Server.UDPConn.WriteToUDP(b, s.Addr); err != nil {
			return err
		}
	}
	return nil
}

// Shutdown
func (s *UDPServer) Shutdown() {
	s.UDPConn.Close()
}

// HandlePacket sends data to consumer
func (s *UDPServer) HandlePacket(p *UDPPacket) error {
	if p.Address != "0" {
		uaddr, err := net.ResolveUDPAddr("udp", p.Address)
		if err != nil {
			return err
		}
		if _, err := s.UDPConn.WriteToUDP(p.Data, uaddr); err != nil {
			return nil
		}
	}
	if err := s.UDPConn.SetDeadline(time.Now().Add(time.Duration(10) * time.Second)); err != nil {
		return err
	}
	return nil
}
