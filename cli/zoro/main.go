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
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/caddyserver/certmagic"
	"github.com/denisbrodbeck/machineid"
	"github.com/go-acme/lego/providers/dns/gcloud"
	"github.com/txthinking/zoro"
	"github.com/txthinking/zoro/https"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "zoro"
	app.Version = "20211230"
	app.Usage = "Expose local TCP and UDP server to external network"
	app.Commands = []*cli.Command{
		&cli.Command{
			Name:  "server",
			Usage: "Run as server mode",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "listen",
					Aliases: []string{"l"},
					Usage:   "Listen address, like: ':9999'",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "Password",
				},
				&cli.StringSliceFlag{
					Name:    "portpassword",
					Aliases: []string{"P"},
					Usage:   "Only allow this port and password, like '1000 password', repeated. If you specify this parameter, --password will be ignored",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("listen") == "" || (c.String("password") == "" && len(c.StringSlice("portpassword")) == 0) {
					cli.ShowCommandHelp(c, "server")
					return nil
				}
				s, err := zoro.NewServer(c.String("listen"), c.String("password"), c.StringSlice("portpassword"))
				if err != nil {
					return err
				}
				go func() {
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					s.Shutdown()
				}()
				return s.ListenAndServe()
			},
		},
		&cli.Command{
			Name:  "client",
			Usage: "Run as client mode",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "server",
					Aliases: []string{"s"},
					Usage:   "Server address, like: 1.2.3.4:9999",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "Password",
				},
				&cli.Int64Flag{
					Name:    "serverport",
					Aliases: []string{"P"},
					Usage:   "Server port you want to use",
				},
				&cli.StringFlag{
					Name:    "client",
					Aliases: []string{"c"},
					Usage:   "Client TCP and/or UDP server address, like: 127.0.0.1:8888",
				},
				&cli.StringFlag{
					Name:    "dir",
					Aliases: []string{"d"},
					Usage:   "Client directory, like: /path/to/www. If you specify this parameter, --client will be ignored",
				},
				&cli.Int64Flag{
					Name:  "dirport",
					Usage: "Work with --dir",
					Value: 8080,
				},
				&cli.Int64Flag{
					Name:  "tcpTimeout",
					Value: 60,
					Usage: "connection tcp keepalive timeout (s)",
				},
				&cli.Int64Flag{
					Name:  "tcpDeadline",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				&cli.Int64Flag{
					Name:  "udpDeadline",
					Value: 60,
					Usage: "connection deadline time (s)",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("server") == "" || c.String("password") == "" || c.Int64("serverport") == 0 {
					cli.ShowCommandHelp(c, "client")
					return nil
				}
				if c.String("client") == "" && (c.String("dir") == "" || c.Int64("dirport") == 0) {
					cli.ShowCommandHelp(c, "client")
					return nil
				}
				cs := c.String("client")
				if c.String("dir") != "" {
					go func() {
						log.Println(http.ListenAndServe(":"+strconv.FormatInt(c.Int64("dirport"), 10), http.FileServer(http.Dir(c.String("dir")))))
					}()
					cs = "127.0.0.1:" + strconv.FormatInt(c.Int64("dirport"), 10)
				}
				s := zoro.NewClient(c.String("server"), c.String("password"), c.Int64("serverport"), "", cs, c.Int64("tcpTimeout"), c.Int64("tcpDeadline"), c.Int64("udpDeadline"))
				go func() {
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					s.Shutdown()
				}()
				return s.ListenAndServe()
			},
		},
		&cli.Command{
			Name:  "httpsserver",
			Usage: "Run as https server mode",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "listen",
					Aliases: []string{"l"},
					Usage:   "Listen address, like: :9999",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "Password",
				},
				&cli.StringFlag{
					Name:    "domain",
					Aliases: []string{"d"},
					Usage:   "Domain, like: domain.com",
				},
				&cli.StringFlag{
					Name:  "cert",
					Usage: "Cert of *.domain.com, like: ./path/to/cert.pem",
				},
				&cli.StringFlag{
					Name:  "key",
					Usage: "Cert key of *.domain.com, like: ./path/to/cert_key.pem",
				},
				&cli.StringFlag{
					Name:  "googledns",
					Usage: "Pointing to a service account file, this will ignore --cert and --key",
				},
				&cli.Int64Flag{
					Name:  "tlsport",
					Usage: "TLS Port, works with --domain",
					Value: 443,
				},
				&cli.Int64Flag{
					Name:  "tlsTimeout",
					Usage: "TLS timeout, works with --domain",
					Value: 60,
				},
				&cli.Int64Flag{
					Name:  "tlsDeadline",
					Usage: "TLS deadline, works with --domain",
					Value: 0,
				},
				&cli.StringSliceFlag{
					Name:    "subdomainpassword",
					Aliases: []string{"P"},
					Usage:   "Only allow this domain and password, like 'subdomain password'. If you specify this parameter, --password will be ignored",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("listen") == "" || c.String("domain") == "" || (c.String("password") == "" && len(c.StringSlice("subdomainpassword")) == 0) {
					cli.ShowCommandHelp(c, "httpsserver")
					return nil
				}
				if (c.String("cert") == "" || c.String("key") == "") && c.String("googledns") == "" {
					cli.ShowCommandHelp(c, "httpsserver")
					return nil
				}
				s, err := https.NewHTTPSServer(c.String("listen"), c.String("password"), c.String("domain"), c.String("cert"), c.String("key"), c.Int64("tlsport"), c.Int64("tlsTimeout"), c.Int64("tlsDeadline"), c.StringSlice("subdomainpassword"))
				if err != nil {
					return err
				}
				if c.String("cert") == "" || c.String("key") == "" {
					certmagic.DefaultACME.Agreed = true
					certmagic.DefaultACME.Email = "cloud+zoro@txthinking.com"
					certmagic.DefaultACME.CA = certmagic.LetsEncryptProductionCA
					if c.String("googledns") != "" {
						b, err := ioutil.ReadFile(c.String("googledns"))
						if err != nil {
							return err
						}
						dp, err := gcloud.NewDNSProviderServiceAccountKey(b)
						if err != nil {
							return err
						}
						certmagic.DefaultACME.DNSProvider = dp
					}
					tc, err := certmagic.TLS([]string{"*." + c.String("domain")})
					if err != nil {
						return err
					}
					tc.NextProtos = []string{"http/1.1"}
					s.TLSConfig = tc
				}
				if err != nil {
					return err
				}
				go func() {
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					s.Shutdown()
				}()
				return s.ListenAndServe()
			},
		},
		&cli.Command{
			Name:  "httpsclient",
			Usage: "Run as https client mode",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "server",
					Aliases: []string{"s"},
					Usage:   "Server address, like: 1.2.3.4:9999",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "Password",
				},
				&cli.StringFlag{
					Name:  "subdomain",
					Usage: "Server subdomain you want to use, default machineid",
				},
				&cli.StringFlag{
					Name:    "client",
					Aliases: []string{"c"},
					Usage:   "Client http 1.1 server address, like: 127.0.0.1:8888",
				},
				&cli.StringFlag{
					Name:    "dir",
					Aliases: []string{"d"},
					Usage:   "Client directory, like: /path/to/www. If you specify this parameter, --client will be ignored",
				},
				&cli.Int64Flag{
					Name:  "dirport",
					Usage: "Work with --dir",
					Value: 8080,
				},
				&cli.Int64Flag{
					Name:  "tcpTimeout",
					Value: 60,
					Usage: "connection tcp keepalive timeout (s)",
				},
				&cli.Int64Flag{
					Name:  "tcpDeadline",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				&cli.Int64Flag{
					Name:  "udpDeadline",
					Value: 60,
					Usage: "connection deadline time (s)",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("server") == "" || c.String("password") == "" {
					cli.ShowCommandHelp(c, "httpsclient")
					return nil
				}
				sd := c.String("subdomain")
				if sd == "" {
					id, err := machineid.ID()
					if err != nil {
						return err
					}
					sd = strings.ReplaceAll(strings.ToLower(id), "-", "")
				}
				fmt.Println("Subdomain:", sd)
				if c.String("client") == "" && (c.String("dir") == "" || c.Int64("dirport") == 0) {
					cli.ShowCommandHelp(c, "httpsclient")
					return nil
				}
				cs := c.String("client")
				if c.String("dir") != "" {
					go func() {
						log.Println(http.ListenAndServe(":"+strconv.FormatInt(c.Int64("dirport"), 10), http.FileServer(http.Dir(c.String("dir")))))
					}()
					cs = "127.0.0.1:" + strconv.FormatInt(c.Int64("dirport"), 10)
				}
				s := zoro.NewClient(c.String("server"), c.String("password"), 0, sd, cs, c.Int64("tcpTimeout"), c.Int64("tcpDeadline"), c.Int64("udpDeadline"))
				go func() {
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					s.Shutdown()
				}()
				return s.ListenAndServe()
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Println(err)
	}
}
