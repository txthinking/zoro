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

	cache "github.com/patrickmn/go-cache"
	"github.com/txthinking/x"
)

// Client .
type Client struct {
	Server       string
	ServerPort   int64
	ClientServer string
	TCPDeadline  int64
	TCPTimeout   int64
	UDPDeadline  int64
	UDPConn      *net.UDPConn
	Cache        *cache.Cache
	Ckv          *x.CryptKV
}

// NewClient .
func NewClient(server, password string, serverPort int64, clientServer string, tcpTimeout, tcpDeadline, udpDeadline int64) *Client {
	c := &Client{
		Server:       server,
		ServerPort:   serverPort,
		ClientServer: clientServer,
		TCPTimeout:   tcpTimeout,
		TCPDeadline:  tcpDeadline,
		UDPDeadline:  udpDeadline,
		Cache:        cache.New(cache.NoExpiration, cache.NoExpiration),
		Ckv: &x.CryptKV{
			AESKey: []byte(password),
		},
	}
	return c
}

// Run .
func (c *Client) Run() error {
	t, err := NewTCPClient(c)
	if err != nil {
		return err
	}
	u, err := NewUDPClient(c)
	if err != nil {
		return err
	}
	errch := make(chan error)
	go func() {
		errch <- t.Run()
	}()
	go func() {
		errch <- u.Run()
	}()
	return <-errch
}
