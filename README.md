# zoro

[中文](README_ZH.md)

zoro (mr2) can help you expose local server to external network. **Support both TCP/UDP**, of course support HTTP/HTTPS. Keep it **simple**, **stupid**.

❤️ A project by [txthinking.com](https://www.txthinking.com)

### Install via [nami](https://github.com/txthinking/nami)

```
nami install zoro
```

### Install via brew

```
brew install zoro
```

### Usage

```
NAME:
   zoro - Expose local TCP and UDP server to external network

USAGE:
   zoro [global options] command [command options] [arguments...]

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
$ zoro server --listen :9999 --password password
```

> More parameters: $ zoro server --help

On local. Assume your remote zoro server is `1.2.3.4:9999`, your local server is `127.0.0.1:8080`, want the remote server to open port `8888`

```
$ zoro client --server 1.2.3.4:9999 --password password --serverport 8888 --client 127.0.0.1:8080
```

> More parameters: $ zoro client --help

Then access `1.2.3.4:8888` equals to access `127.0.0.1:8080`

## `httpsserver` and `httpsclient`

On remote server. Assume your domain is `domain.com`, cert of `*.domain.com` is `./domain_com_cert.pem` and `./domain_com_cert_key.pem`, want https listen on `443`. Note that the firewall opens TCP on all relevant ports

```
$ zoro httpsserver --listen :9999 --password password --domain domain.com --cert ./domain_com_cert.pem --key ./domain_com_cert_key.pem --tlsport 443
```

> More parameters: $ zoro httpsserver --help

On local. Assume your remote zoro httpsserver is `1.2.3.4:9999`, your local HTTP 1.1 server is `127.0.0.1:8080`, want the remote server to open subdomain `hello`

```
$ zoro httpsclient --server 1.2.3.4:9999 --password password --subdomain hello --client 127.0.0.1:8080
```

> More parameters: $ zoro httpsclient --help

Then access `https://hello.domain.com:443` equals to access `http://127.0.0.1:8080`

## Example of `server` and `client`

#### Expose local HTTP server

```
$ zoro client --server 1.2.3.4:9999 --password password --serverport 8888 --client 127.0.0.1:8080
```

Then access `1.2.3.4:8888` equals to access `127.0.0.1:8080`

#### Expose local SSH

```
$ zoro client --server 1.2.3.4:9999 --password password --serverport 8888 --client 127.0.0.1:22
```

Then access `1.2.3.4:8888` equals to access `127.0.0.1:22`

```
$ ssh -oPort=8888 yourlocaluser@1.2.3.4
```

#### Expose local DNS server

```
$ zoro client --server 1.2.3.4:9999 --password password --serveport 8888 --client 127.0.0.1:53
```

Then access `1.2.3.4:8888` equals to access `127.0.0.1:53`

```
$ dig github.com @1.2.3.4 -p 8888
```

#### Expose local directory via HTTP

```
$ zoro client --server 1.2.3.4:9999 --password password --serverport 8888 --dir /path/to/www --dirport 8080
```

Then access `1.2.3.4:8888` equals to access `127.0.0.1:8080`, web root is /path/to/www

#### Expose any TCP/UDP service

```
...
```

## About UDP

In some cases of multi-layer NAT, UDP may fail. I passed the test when I connected directly to the Wi-Fi provided by the ISP.

## License

Licensed under The GPLv3 License
