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
	"net"

	"github.com/txthinking/crypto"
)

// Client .
type Client struct {
	Server       string
	ServerPort   int64
	ServerDomain string
	ClientServer string
	TCPDeadline  int64
	TCPTimeout   int64
	UDPDeadline  int64
	UDPConn      *net.UDPConn
	Ckv          *crypto.KV
	TCPClient    *TCPClient
	UDPClient    *UDPClient
}

// NewClient .
func NewClient(server, password string, serverPort int64, serverDomain, clientServer string, tcpTimeout, tcpDeadline, udpDeadline int64) *Client {
	c := &Client{
		Server:       server,
		ServerPort:   serverPort,
		ServerDomain: serverDomain,
		ClientServer: clientServer,
		TCPTimeout:   tcpTimeout,
		TCPDeadline:  tcpDeadline,
		UDPDeadline:  udpDeadline,
		Ckv: &crypto.KV{
			AESKey: []byte(password),
		},
	}
	return c
}

// Run .
func (c *Client) ListenAndServe() error {
	var err error
	c.TCPClient, err = NewTCPClient(c)
	if err != nil {
		return err
	}
	if c.ServerDomain != "" {
		return c.TCPClient.Run()
	}
	c.UDPClient, err = NewUDPClient(c)
	if err != nil {
		return err
	}
	errch := make(chan error)
	go func() {
		errch <- c.TCPClient.Run()
	}()
	go func() {
		errch <- c.UDPClient.Run()
	}()
	return <-errch
}

// Shutdown server.
func (c *Client) Shutdown() error {
	if c.TCPClient != nil {
		c.TCPClient.Stop()
	}
	if c.UDPClient != nil {
		c.UDPClient.Stop()
	}
	return nil
}
