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
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"

	"github.com/urfave/cli"
)

var debug bool
var debugListen string

func main() {
	app := cli.NewApp()
	app.Name = "Mr.2"
	app.Version = "20190501"
	app.Usage = "Expose local server to external network"
	app.Author = "Cloud"
	app.Email = "cloud@txthinking.com"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "Enable debug, more logs",
			Destination: &debug,
		},
		cli.StringFlag{
			Name:        "listen, l",
			Usage:       "Listen address for debug",
			Value:       "127.0.0.1:6060",
			Destination: &debugListen,
		},
	}
	app.Commands = []cli.Command{
		cli.Command{
			Name:  "server",
			Usage: "Run as server mode",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "listen, l",
					Usage: "Listen address, like: 1.2.3.4:5",
				},
				cli.StringFlag{
					Name:  "password, p",
					Usage: "Password",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("listen") == "" || c.String("password") == "" {
					cli.ShowCommandHelp(c, "server")
					return nil
				}
				if debug {
					go func() {
						log.Println(http.ListenAndServe(debugListen, nil))
					}()
				}
				return RunServer(c.String("listen"), c.String("password"))
			},
		},
		cli.Command{
			Name:  "client",
			Usage: "Run as client mode",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "server, s",
					Usage: "Server address, like: 1.2.3.4:5",
				},
				cli.StringFlag{
					Name:  "password, p",
					Usage: "Password",
				},
				cli.Int64Flag{
					Name:  "serverPort, P",
					Usage: "Server port you want to use",
				},
				cli.StringFlag{
					Name:  "clientServer, c",
					Usage: "Client server address, like: 1.2.3.4:5",
				},
				cli.StringFlag{
					Name:  "clientDirectory",
					Usage: "Client directory, like: /path/to/www",
				},
				cli.Int64Flag{
					Name:  "clientPort",
					Usage: "With clientDirectory",
					Value: 54321,
				},
				cli.Int64Flag{
					Name:  "tcpTimeout",
					Value: 60,
					Usage: "connection tcp keepalive timeout (s)",
				},
				cli.Int64Flag{
					Name:  "tcpDeadline",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				cli.Int64Flag{
					Name:  "udpDeadline",
					Value: 60,
					Usage: "connection deadline time (s)",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("server") == "" || c.String("password") == "" || c.Int64("serverPort") == 0 {
					cli.ShowCommandHelp(c, "client")
					return nil
				}
				if c.String("clientServer") == "" && c.String("clientDirectory") == "" {
					cli.ShowCommandHelp(c, "client")
					return nil
				}
				if debug {
					go func() {
						log.Println(http.ListenAndServe(debugListen, nil))
					}()
				}
				cs := c.String("clientServer")
				if c.String("clientDirectory") != "" {
					go func() {
						log.Println(http.ListenAndServe(":"+strconv.FormatInt(c.Int64("clientPort"), 10), http.FileServer(http.Dir(c.String("clientDirectory")))))
					}()
					cs = "localhost:" + strconv.FormatInt(c.Int64("clientPort"), 10)
				}
				h, _, err := net.SplitHostPort(c.String("server"))
				if err != nil {
					return err
				}
				s := net.JoinHostPort(h, strconv.FormatInt(c.Int64("serverPort"), 10))
				fmt.Println(s)
				return RunClient(c.String("server"), c.String("password"), c.Int64("serverPort"), cs, c.Int64("tcpTimeout"), c.Int64("tcpDeadline"), c.Int64("udpDeadline"))
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
