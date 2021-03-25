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
	"errors"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/txthinking/encrypt"
)

// HTTPSServer .
type HTTPSServer struct {
	TCPAddr   *net.TCPAddr
	TCPListen *net.TCPListener
	Ckv       *encrypt.KV

	Domain      string
	TLSConfig   *tls.Config
	TLSAddr     *net.TCPAddr
	TLSTimeout  int64
	TLSDeadline int64
	DomainCkv   map[string]*encrypt.KV
	TLSServer   *TLSServer
}

func NewHTTPSServer(addr, password string, domain, cert, certKey string, tlsPort, tlsTimeout, tlsDeadline int64, domainPassword []string) (*HTTPSServer, error) {
	if tlsPort == 0 {
		return nil, errors.New("Your forgot tlsPort")
	}
	tlsAddr, err := net.ResolveTCPAddr("tcp", ":"+strconv.FormatInt(tlsPort, 10))
	if err != nil {
		return nil, err
	}
	cer, err := tls.LoadX509KeyPair(cert, certKey)
	if err != nil {
		return nil, err
	}
	tc := &tls.Config{Certificates: []tls.Certificate{cer}}
	tc.NextProtos = []string{"http/1.1"}
	dc := make(map[string]*encrypt.KV)
	for _, v := range domainPassword {
		l := strings.Split(v, " ")
		if len(l) != 2 {
			return nil, errors.New("Wrong format: " + v)
		}
		ckv := &encrypt.KV{
			AESKey: []byte(l[1]),
		}
		dc[l[0]] = ckv
	}
	taddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	s := &HTTPSServer{
		TCPAddr: taddr,
		Ckv: &encrypt.KV{
			AESKey: []byte(password),
		},
		Domain:      domain,
		TLSConfig:   tc,
		TLSAddr:     tlsAddr,
		TLSTimeout:  tlsTimeout,
		TLSDeadline: tlsDeadline,
		DomainCkv:   dc,
	}
	return s, nil
}

func (s *HTTPSServer) ListenAndServe() error {
	return s.RunTLSServer()
}

func (s *HTTPSServer) RunTLSServer() error {
	var err error
	s.TLSServer, err = NewTLSServer(s)
	if err != nil {
		return err
	}
	defer s.TLSServer.Shutdown()
	errch := make(chan error)
	go func() {
		if err := s.TLSServer.ListenAndServe(); err != nil {
			errch <- err
		}
	}()
	go func() {
		s.TCPListen, err = net.ListenTCP("tcp", s.TCPAddr)
		if err != nil {
			errch <- err
			return
		}
		defer s.TCPListen.Close()
		for {
			c, err := s.TCPListen.AcceptTCP()
			if err != nil {
				errch <- err
				return
			}
			go func(c *net.TCPConn) {
				defer c.Close()
				if err := s.TLSServer.HandleClient(c); err != nil {
					log.Println(err)
				}
			}(c)
		}
	}()
	return <-errch
}

func (s *HTTPSServer) Shutdown() {
	s.TCPListen.Close()
}
