# mr2

[English](README.md)

[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](http://www.gnu.org/licenses/gpl-3.0)
[![捐赠](https://img.shields.io/badge/%E6%94%AF%E6%8C%81-%E6%8D%90%E8%B5%A0-ff69b4.svg)](https://github.com/sponsors/txthinking)
[![交流群](https://img.shields.io/badge/%E7%94%B3%E8%AF%B7%E5%8A%A0%E5%85%A5-%E4%BA%A4%E6%B5%81%E7%BE%A4-ff69b4.svg)](https://docs.google.com/forms/d/e/1FAIpQLSdzMwPtDue3QoezXSKfhW88BXp57wkbDXnLaqokJqLeSWP9vQ/viewform)

mr2 帮助你将本地端口暴露在外网.**支持TCP/UDP**, 当然也支持HTTP. Keep it **simple**, **stupid**.

❤️ A project by [txthinking.com](https://www.txthinking.com)

### 使用[nami](https://github.com/txthinking/nami)安装

```
$ nami install github.com/txthinking/mr2
```

### 使用brew安装（macOS）

```
$ brew install mr2
```

### 使用说明

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

## `服务端` 及 `客户端` 使用教程

在远程服务器上. 注意防火墙开放所有相关端口的TCP和UDP协议

```
$ mr2 server -l :9999 -p password
```

> 更多参数: $ mr2 server -h

在本地. 假设你的远程 mr2 server 是`1.2.3.4:9999`, 你的本地服务是`127.0.0.1:8080`, 你想让远程服务器开放`8888`端口

```
$ mr2 client -s 1.2.3.4:9999 -p password --serverPort 8888 -c 127.0.0.1:8080
```

> 更多参数: $ mr2 client -h<br/>

现在访问 `1.2.3.4:8888` 就等于 `127.0.0.1:8080`

## `服务端` 及 `客户端` 的使用例子

#### 暴露本地HTTP服务

```
$ mr2 client -s 1.2.3.4:9999 -p password --serverPort 8888 -c 127.0.0.1:8080
```

现在访问 `1.2.3.4:8888` 就等于 `127.0.0.1:8080`

#### 暴露本地SSH服务

```
$ mr2 client -s 1.2.3.4:9999 -p password --serverPort 8888 -c 127.0.0.1:22
```

现在访问 `1.2.3.4:8888` 就等于 `127.0.0.1:22`

```
$ ssh -oPort=8888 yourlocaluser@1.2.3.4
```

#### 暴露本地DNS服务

```
$ mr2 client -s 1.2.3.4:9999 -p password --serverPort 8888 -c 127.0.0.1:53
```

现在访问 `1.2.3.4:8888` 就等于 `127.0.0.1:53`

```
$ dig github.com @1.2.3.4 -p 8888
```

#### 暴露本地目录通过HTTP

```
$ mr2 client -s 1.2.3.4:9999 -p password --serverPort 8888 --clientDirectory /path/to/www --clientPort 8080
```

现在访问 `1.2.3.4:8888` 就等于 `127.0.0.1:8080`, web root 是 /path/to/www

#### 暴露本地brook代理

```
$ brook server -l :8080 -p password # 或者用 wsserver
$ mr2 client -s 1.2.3.4:9999 -p password --serverPort 8888 -c 127.0.0.1:8080
```

现在访问 `1.2.3.4:8888` 就等于访问 `127.0.0.1:8080`，做到在没有公网ip的机器上搭建brook，在外使用家里或者公司非公网ip上网

#### 暴露你能想到的任何TCP/UDP服务

```
...
```

## `httpsserver` 以及 `httpsclient`

在远程服务器上. 假设你的域名是 `domain.com`, 泛域名证书`*.domain.com` 是 `./domain_com_cert.pem` 和 `./domain_com_cert_key.pem`, 想让HTTPS监听 443`. 注意防火墙开放任何相关端口的TCP协议

```
$ mr2 httpsserver -l :9999 -p password --domain domain.com --cert ./domain_com_cert.pem --certKey ./domain_com_cert_key.pem --tlsPort 443
```

> 更多参数: $ mr2 httpsserver -h

在本地. 假设你的远程 mr2 httpsserver 是 `1.2.3.4:9999`, 你的本地 HTTP 1.1 服务是 `127.0.0.1:8080`, 想让远程服务器开放子域名 `hey`

```
$ mr2 httpsclient -s 1.2.3.4:9999 -p password --serverSubdomain hey -c 127.0.0.1:8080
```

> 更多参数: $ mr2 httpsclient -h

现在访问 `https://hey.domain.com:443` 就等于 `http://127.0.0.1:8080`

## 关于UDP

在一些多层NAT情况下, 可能UDP会失败. 我在本地直接连接ISP提供的Wi-Fi的情况测试通过.

## 开源协议

基于 GPLv3 协议开源
