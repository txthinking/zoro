# Mr.2

[![Build Status](https://travis-ci.org/txthinking/mr2.svg?branch=master)](https://travis-ci.org/txthinking/mr2) [![Go Report Card](https://goreportcard.com/badge/github.com/txthinking/mr2)](https://goreportcard.com/report/github.com/txthinking/mr2) [![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](http://www.gnu.org/licenses/gpl-3.0) [![Wiki](https://img.shields.io/badge/docs-wiki-blue.svg)](https://github.com/txthinking/mr2/wiki)

---

Coming soon...

---

### Table of Contents

* [What is Mr.2](#what-is-mr2)
* [Download](#download)
* [**Server**](#server)
* [**Client**](#client)
* [Contributing](#contributing)
* [License](#license)

## What is Mr.2

Mr.2 can help you expose local server to external network.<br/>
Mr.2 support both TCP/UDP, of course support HTTP.<br/>
Keep it **simple**, **stupid**.

## Download

| Download | OS | Arch |
| --- | --- | --- |
| [mr2](https://github.com/txthinking/mr2/releases/download/v20190501/mr2) | Linux | amd64 |
| [mr2_darwin_amd64](https://github.com/txthinking/mr2/releases/download/v20190501/mr2_darwin_amd64) | MacOS | amd64 |
| [mr2_windows_amd64.exe](https://github.com/txthinking/mr2/releases/download/v20190501/mr2_windows_amd64.exe) | Windows | amd64 |

**See [releases](https://github.com/txthinking/mr2/releases) for other platforms**

## Mr.2

```
NAME:
   Mr.2 - Expose local server to external network

USAGE:
   mr2 [global options] command [command options] [arguments...]

VERSION:
   20190501

AUTHOR:
   Cloud <cloud@txthinking.com>

COMMANDS:
     server   Run as server mode
     client   Run as client mode
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug, -d               Enable debug, more logs
   --listen value, -l value  Listen address for debug (default: "127.0.0.1:6060")
   --help, -h                show help
   --version, -v             print the version
```

### Server

```
$ mr2 server -l :9999 -p password
```

### Client

```
# Suppose your local server address is 127.0.0.1:1234

$ mr2 client -s server_address:port -p password -P 5678 -c 127.0.0.1:1234

# Now the external server address is: server_address:5678
```

```
# Suppose your local web root directory is /path/to/www

$ mr2 client -s server_address:port -p password -P 5678 --clientDiretory /path/to/www

# Now the external server address is: server_address:5678
```

## Contributing

Please read [CONTRIBUTING.md](https://github.com/txthinking/mr2/blob/master/.github/CONTRIBUTING.md) first

## License

Licensed under The GPLv3 License
