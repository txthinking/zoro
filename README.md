# Mr.2

[![Build Status](https://travis-ci.org/txthinking/mr2.svg?branch=master)](https://travis-ci.org/txthinking/mr2) [![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](http://www.gnu.org/licenses/gpl-3.0)
[![ZH](https://img.shields.io/badge/%E4%B8%AD%E6%96%87-README-blue.svg)](https://github.com/txthinking/mr2/blob/master/README_zh.md)
[![Financial Contributors on Open Collective](https://opencollective.com/txthinking-mr2/all/badge.svg?label=financial+contributors)](https://opencollective.com/txthinking-mr2) 

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
| [mr2](https://github.com/txthinking/mr2/releases/download/v20190616/mr2) | Linux | amd64 |
| [mr2_darwin_amd64](https://github.com/txthinking/mr2/releases/download/v20190616/mr2_darwin_amd64) | MacOS | amd64 |
| [mr2_windows_amd64.exe](https://github.com/txthinking/mr2/releases/download/v20190616/mr2_windows_amd64.exe) | Windows | amd64 |

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

## Contributors

### Code Contributors

This project exists thanks to all the people who contribute. [[Contribute](CONTRIBUTING.md)].
<a href="https://github.com/txthinking/mr2/graphs/contributors"><img src="https://opencollective.com/txthinking-mr2/contributors.svg?width=890&button=false" /></a>

### Financial Contributors

Become a financial contributor and help us sustain our community. [[Contribute](https://opencollective.com/txthinking-mr2/contribute)]

#### Individuals

<a href="https://opencollective.com/txthinking-mr2"><img src="https://opencollective.com/txthinking-mr2/individuals.svg?width=890"></a>

#### Organizations

Support this project with your organization. Your logo will show up here with a link to your website. [[Contribute](https://opencollective.com/txthinking-mr2/contribute)]

<a href="https://opencollective.com/txthinking-mr2/organization/0/website"><img src="https://opencollective.com/txthinking-mr2/organization/0/avatar.svg"></a>
<a href="https://opencollective.com/txthinking-mr2/organization/1/website"><img src="https://opencollective.com/txthinking-mr2/organization/1/avatar.svg"></a>
<a href="https://opencollective.com/txthinking-mr2/organization/2/website"><img src="https://opencollective.com/txthinking-mr2/organization/2/avatar.svg"></a>
<a href="https://opencollective.com/txthinking-mr2/organization/3/website"><img src="https://opencollective.com/txthinking-mr2/organization/3/avatar.svg"></a>
<a href="https://opencollective.com/txthinking-mr2/organization/4/website"><img src="https://opencollective.com/txthinking-mr2/organization/4/avatar.svg"></a>
<a href="https://opencollective.com/txthinking-mr2/organization/5/website"><img src="https://opencollective.com/txthinking-mr2/organization/5/avatar.svg"></a>
<a href="https://opencollective.com/txthinking-mr2/organization/6/website"><img src="https://opencollective.com/txthinking-mr2/organization/6/avatar.svg"></a>
<a href="https://opencollective.com/txthinking-mr2/organization/7/website"><img src="https://opencollective.com/txthinking-mr2/organization/7/avatar.svg"></a>
<a href="https://opencollective.com/txthinking-mr2/organization/8/website"><img src="https://opencollective.com/txthinking-mr2/organization/8/avatar.svg"></a>
<a href="https://opencollective.com/txthinking-mr2/organization/9/website"><img src="https://opencollective.com/txthinking-mr2/organization/9/avatar.svg"></a>

## License

Licensed under The GPLv3 License

