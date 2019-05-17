# Mr.2

[![Build Status](https://travis-ci.org/txthinking/mr2.svg?branch=master)](https://travis-ci.org/txthinking/mr2) [![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](http://www.gnu.org/licenses/gpl-3.0)
[![ZH](https://img.shields.io/badge/%E4%B8%AD%E6%96%87-README-blue.svg)](https://github.com/txthinking/mr2/blob/master/README_zh.md)

### Table of Contents

* [What is Mr.2](#what-is-mr2)
* [Download](#download)
* [**Server**](#server)
* [**Client**](#client)
* [Example](#example)
* [Contributing](#contributing)
* [License](#license)

## What is Mr.2

Mr.2 can help you expose local server to external network. Support both TCP/UDP, of course support HTTP.<br/>
Keep it **simple**, **stupid**.

## Download

| Download | OS | Arch |
| --- | --- | --- |
| [mr2](https://github.com/txthinking/mr2/releases/download/v20190518/mr2) | Linux | amd64 |
| [mr2_darwin_amd64](https://github.com/txthinking/mr2/releases/download/v20190518/mr2_darwin_amd64) | MacOS | amd64 |
| [mr2_windows_amd64.exe](https://github.com/txthinking/mr2/releases/download/v20190518/mr2_windows_amd64.exe) | Windows | amd64 |

See [releases](https://github.com/txthinking/mr2/releases) for other platforms. Or `go get github.com/txthinking/mr2/cli/mr2`.

### Server

```
$ mr2 server -l :9999 -p password
```

```
# Only allow partial ports, and set password on each port
$ mr2 server -l :9999 -P '5678 password' -P '6789 password1'
```

### Client

```
# Local server is 127.0.0.1:1234, expect to expose: server_address:5678
$ mr2 client -s server_address:port -p password -P 5678 -c 127.0.0.1:1234
```

```
# Local web root is /path/to/www, expect to expose: server_address:5678
$ mr2 client -s server_address:port -p password -P 5678 --clientDirectory /path/to/www
```

### Example

#### Access local HTTP server

```
$ mr2 client -s server_address:port -p password -P 5678 -c 127.0.0.1:8080

# then
Your HTTP server in external network is: server_address:5678
```

#### SSH into local computer

```
$ mr2 client -s server_address:port -p password -P 5678 -c 127.0.0.1:22

# then
$ ssh -oPort=5678 user@server_address
```

#### Access local DNS server

```
$ mr2 client -s server_address:port -p password -P 5678 -c 127.0.0.1:53

# then
Your DNS server in external network is: server_address:5678

$ dig github.com @server_address -p 5678
```

#### Access your local directory via HTTP

```
$ mr2 client -s server_address:port -p password -P 5678 --clientDirectory /path/to/www

# then
A HTTP server in external network is: server_address:5678
```

#### Any TCP-based/UDP-based ideas you think of

...

## Contributing

Please read [CONTRIBUTING.md](https://github.com/txthinking/mr2/blob/master/.github/CONTRIBUTING.md) first

## License

Licensed under The GPLv3 License
