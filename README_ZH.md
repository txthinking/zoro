# Mr.2

[English](README.md)

[![Build Status](https://travis-ci.org/txthinking/mr2.svg?branch=master)](https://travis-ci.org/txthinking/mr2) [![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](http://www.gnu.org/licenses/gpl-3.0)
[![捐赠](https://img.shields.io/badge/%E6%94%AF%E6%8C%81-%E6%8D%90%E8%B5%A0-ff69b4.svg)](https://www.txthinking.com/opensource-support.html)
[![交流群](https://img.shields.io/badge/%E7%94%B3%E8%AF%B7%E5%8A%A0%E5%85%A5-%E4%BA%A4%E6%B5%81%E7%BE%A4-ff69b4.svg)](https://docs.google.com/forms/d/e/1FAIpQLSdzMwPtDue3QoezXSKfhW88BXp57wkbDXnLaqokJqLeSWP9vQ/viewform)

## 什么是 Mr.2

Mr.2 帮助你将本地端口暴露在外网.支持TCP/UDP, 当然也支持HTTP. Keep it **simple**, **stupid**.

### 用 [nami](https://github.com/txthinking/nami) 安装

```
$ nami install github.com/txthinking/mr2
```

或直接下载二进制命令文件 [releases](https://github.com/txthinking/mr2/releases)

### Server

    $ mr2 server -l :9999 -p password

    # Only allow partial ports, and set password on each port
    $ mr2 server -l :9999 -P '5678 password' -P '6789 password1'

### Client

    # Local server is 127.0.0.1:1234, expect to expose: server_address:5678
    $ mr2 client -s server_address:port -p password -P 5678 -c 127.0.0.1:1234

    # Local web root is /path/to/www, expect to expose: server_address:5678
    $ mr2 client -s server_address:port -p password -P 5678 --clientDirectory /path/to/www

### 举例

#### Access local HTTP server

    $ mr2 client -s server_address:port -p password -P 5678 -c 127.0.0.1:8080

    # then
    Your HTTP server in external network is: server_address:5678

#### SSH into local computer

    $ mr2 client -s server_address:port -p password -P 5678 -c 127.0.0.1:22

    # then
    $ ssh -oPort=5678 user@server_address

#### Access local DNS server

    $ mr2 client -s server_address:port -p password -P 5678 -c 127.0.0.1:53

    # then
    Your DNS server in external network is: server_address:5678

    $ dig github.com @server_address -p 5678

#### Access your local directory via HTTP

    $ mr2 client -s server_address:port -p password -P 5678 --clientDirectory /path/to/www

    # then
    A HTTP server in external network is: server_address:5678

#### Any TCP-based/UDP-based ideas you think of

...

## 贡献

请先阅读 [CONTRIBUTING.md](https://github.com/txthinking/mr2/blob/master/.github/CONTRIBUTING.md)

## 开源协议

基于 GPLv3 协议开源
