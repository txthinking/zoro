# mr2

[中文](README_ZH.md)

[![Build Status](https://travis-ci.org/txthinking/mr2.svg?branch=master)](https://travis-ci.org/txthinking/mr2) [![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](http://www.gnu.org/licenses/gpl-3.0)
[![Donate](https://img.shields.io/badge/Support-Donate-ff69b4.svg)](https://www.txthinking.com/opensource-support.html)
[![Slack](https://img.shields.io/badge/Join-Slack-ff69b4.svg)](https://docs.google.com/forms/d/e/1FAIpQLSdzMwPtDue3QoezXSKfhW88BXp57wkbDXnLaqokJqLeSWP9vQ/viewform)

mr2 can help you expose local server to external network. **Support both TCP/UDP**, of course support HTTP. Keep it **simple**, **stupid**.

❤️ A project by [txthinking.com](https://www.txthinking.com)

### Install via [nami](https://github.com/txthinking/nami)

```
$ nami install github.com/txthinking/mr2
```

### Install via brew (macOS only)

```
$ brew install mr2
```

### Usage

```
NAME:
   mr2 - Expose local TCP and UDP server to external network

USAGE:
   mr2 [global options] command [command options] [arguments...]

VERSION:
   20210401

COMMANDS:
   server       Run as server mode
   client       Run as client mode
   httpsserver  Run as https server mode
   httpsclient  Run as https client mode
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

## `server` and `client`

On remote server. Note that the firewall opens TCP and UDP on all relevant ports

```
$ mr2 server -l :9999 -p password
```

> More parameters: $ mr2 server -h

On local. Assume your remote mr2 server is `1.2.3.4:9999`, your local server is `127.0.0.1:8080`, want the remote server to open port `8888`

```
$ mr2 client -s 1.2.3.4:9999 -p password --serverPort 8888 -c 127.0.0.1:8080
```

> More parameters: $ mr2 client -h<br/>

Then access `1.2.3.4:8888` equals to access `127.0.0.1:8080`

## Example of `server` and `client` 

#### Expose local HTTP server

```
$ mr2 client -s 1.2.3.4:9999 -p password --serverPort 8888 -c 127.0.0.1:8080
```

Then access `1.2.3.4:8888` equals to access `127.0.0.1:8080`

#### Expose local SSH

```
$ mr2 client -s 1.2.3.4:9999 -p password --serverPort 8888 -c 127.0.0.1:22
```

Then access `1.2.3.4:8888` equals to access `127.0.0.1:22`

```
$ ssh -oPort=8888 yourlocaluser@1.2.3.4
```

#### Expose local DNS server

```
$ mr2 client -s 1.2.3.4:9999 -p password --serverPort 8888 -c 127.0.0.1:53
```

Then access `1.2.3.4:8888` equals to access `127.0.0.1:53`

```
$ dig github.com @1.2.3.4 -p 8888
```

#### Expose local directory via HTTP

```
$ mr2 client -s 1.2.3.4:9999 -p password --serverPort 8888 --clientDirectory /path/to/www --clientPort 8080
```

Then access `1.2.3.4:8888` equals to access `127.0.0.1:8080`, web root is /path/to/www

#### Expose local brook server

```
$ brook server -l :8080 -p password # or wsserver
$ mr2 client -s 1.2.3.4:9999 -p password --serverPort 8888 -c 127.0.0.1:8080
```

Then access `1.2.3.4:8888` equals to access `127.0.0.1:8080`, used to create a brook server or wsserver in a server even if there is no public IP

#### Expose any TCP/UDP service

```
...
```

## `httpsserver` and `httpsclient`

On remote server. Assume your domain is `domain.com`, cert of `*.domain.com` is `./domain_com_cert.pem` and `./domain_com_cert_key.pem`, want https listen on `443`. Note that the firewall opens TCP on all relevant ports

```
$ mr2 httpsserver -l :9999 -p password --domain domain.com --cert ./domain_com_cert.pem --certKey ./domain_com_cert_key.pem --tlsPort 443
```

> More parameters: $ mr2 httpsserver -h<br/>

On local. Assume your remote mr2 httpsserver is `1.2.3.4:9999`, your local HTTP 1.1 server is `127.0.0.1:8080`, want the remote server to open subdomain `hey`

```
$ mr2 httpsclient -s 1.2.3.4:9999 -p password --serverSubdomain hey -c 127.0.0.1:8080
```

> More parameters: $ mr2 httpsclient -h

Then access `https://hey.domain.com:443` equals to access `http://127.0.0.1:8080`

## About UDP

In some cases of multi-layer NAT, UDP may fail. I passed the test when I connected directly to the Wi-Fi provided by the ISP.

## License

Licensed under The GPLv3 License
