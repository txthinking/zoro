# zoro

[English](README.md)

zoro (mr2) 帮助你将本地端口暴露在外网.**支持TCP/UDP**, 当然也支持HTTP/HTTPS. Keep it **simple**, **stupid**.

❤️ A project by [txthinking.com](https://www.txthinking.com)

### 使用[nami](https://github.com/txthinking/nami)安装

```
nami install zoro
```

### 使用brew安装

```
brew install zoro
```

### 使用说明

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

## `server` 及 `client`

在远程服务器上. 注意防火墙开放所有相关端口的TCP和UDP协议

```
$ zoro server --listen :9999 --password password
```

> 更多参数: $ zoro server --help

在本地. 假设你的远程 zoro server 是`1.2.3.4:9999`, 你的本地服务是`127.0.0.1:8080`, 你想让远程服务器开放`8888`端口

```
$ zoro client --server 1.2.3.4:9999 --password password --serverport 8888 --client 127.0.0.1:8080
```

> 更多参数: $ zoro client --help

现在访问 `1.2.3.4:8888` 就等于 `127.0.0.1:8080`

## `httpsserver` 及 `httpsclient`

在远程服务器上. 假设你的域名是 `domain.com`, 泛域名证书`*.domain.com` 是 `./domain_com_cert.pem` 和 `./domain_com_cert_key.pem`, 想让HTTPS监听 `443`. 注意防火墙开放任何相关端口的TCP协议

```
$ zoro httpsserver --listen :9999 --password password --domain domain.com --cert ./domain_com_cert.pem --key ./domain_com_cert_key.pem --tlsport 443
```

> 更多参数: $ zoro httpsserver --help

在本地. 假设你的远程 zoro httpsserver 是 `1.2.3.4:9999`, 你的本地 HTTP 1.1 服务是 `127.0.0.1:8080`, 想让远程服务器开放子域名 `hello`

```
$ zoro httpsclient --server 1.2.3.4:9999 --password password --subdomain hello --client 127.0.0.1:8080
```

> 更多参数: $ zoro httpsclient --help

现在访问 `https://hello.domain.com:443` 就等于 `http://127.0.0.1:8080`

## `server` 及 `client` 的使用例子

#### 暴露本地HTTP服务

```
$ zoro client --server 1.2.3.4:9999 --password password --serverport 8888 --client 127.0.0.1:8080
```

现在访问 `1.2.3.4:8888` 就等于 `127.0.0.1:8080`

#### 暴露本地SSH服务

```
$ zoro client --server 1.2.3.4:9999 --password password --serverport 8888 --client 127.0.0.1:22
```

现在访问 `1.2.3.4:8888` 就等于 `127.0.0.1:22`

```
$ ssh -oPort=8888 yourlocaluser@1.2.3.4
```

#### 暴露本地DNS服务

```
$ zoro client --server 1.2.3.4:9999 --password password --serveport 8888 --client 127.0.0.1:53
```

现在访问 `1.2.3.4:8888` 就等于 `127.0.0.1:53`

```
$ dig github.com @1.2.3.4 -p 8888
```

#### 暴露本地目录通过HTTP

```
$ zoro client --server 1.2.3.4:9999 --password password --serverport 8888 --dir /path/to/www --dirport 8080
```

现在访问 `1.2.3.4:8888` 就等于 `127.0.0.1:8080`, web root 是 /path/to/www

#### 暴露你能想到的任何TCP/UDP服务

```
...
```

## 关于UDP

在一些多层NAT情况下, 可能UDP会失败. 我在本地直接连接ISP提供的Wi-Fi的情况测试通过.

## 开源协议

基于 GPLv3 协议开源
