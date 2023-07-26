# Pingtunnel

[<img src="https://img.shields.io/github/license/esrrhs/pingtunnel">](https://github.com/esrrhs/pingtunnel)
[<img src="https://img.shields.io/github/languages/top/esrrhs/pingtunnel">](https://github.com/esrrhs/pingtunnel)
[![Go Report Card](https://goreportcard.com/badge/github.com/esrrhs/pingtunnel)](https://goreportcard.com/report/github.com/esrrhs/pingtunnel)
[<img src="https://img.shields.io/github/v/release/esrrhs/pingtunnel">](https://github.com/esrrhs/pingtunnel/releases)
[<img src="https://img.shields.io/github/downloads/esrrhs/pingtunnel/total">](https://github.com/esrrhs/pingtunnel/releases)
[<img src="https://img.shields.io/docker/pulls/esrrhs/pingtunnel">](https://hub.docker.com/repository/docker/esrrhs/pingtunnel)
[<img src="https://img.shields.io/github/actions/workflow/status/esrrhs/pingtunnel/go.yml?branch=master">](https://github.com/esrrhs/pingtunnel/actions)

Pingtunnel is a tool that send TCP/UDP traffic over ICMP.

## Note: This tool is only to be used for study and research, do not use it for illegal purposes

![image](network.jpg)

## Usage

```shell
NAME:
   pingtunnel - A tool that send TCP/UDP traffic over ICMP

USAGE:
   pingtunnel [global options] command [command options] [arguments...]

VERSION:
   Git:[DEV] (go1.20)

COMMANDS:
   server   Run a PingTunnel Server
   client   Run a PingTunnel Client: Sock5 Proxy
   tunnel   Run a PingTunnel Tunnel: TCP/UDP Port Forward
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

### Install server

- First prepare a server with a public IP, such as EC2 on AWS, assuming the domain name or public IP is example.com
- Download the corresponding installation package from [releases](https://github.com/esrrhs/pingtunnel/releases), such as pingtunnel_linux64.zip, then decompress and execute with **root** privileges
- `--key` parameter is **INT** type, only supports numbers between `0-2147483647`

```shell
wget https://github.com/honwen/pingtunnel/releases/download/v{version}/pingtunnel-linux-amd64-v{version}.tar.gz # {link of latest release}
tar -zxvf pingtunnel-linux-amd64-*.tar.gz

pingtunnel server
```

```shell
NAME:
   pingtunnel server - Run a PingTunnel Server

USAGE:
   pingtunnel server [command options] [arguments...]

OPTIONS:
   --key key             numeric key in [0-2147483647] (default: 0)
   --maxconn connection  max num of connections, 0 means no limit (default: 0)
   --maxthread thread    max process thread in server (default: 100)
   --maxbuffer buffer    max process thread's buffer in server (default: 1000)
   --timeout ms          timeout(ms) period for the server to initiate a connection to the destination address (default: 1000)
```

- (Optional) Disable system default ping

```shell
sysctl -w net.ipv4.icmp_echo_ignore_all=1
# or
echo 1 > /proc/sys/net/ipv4/icmp_echo_ignore_all
```

### Install the client

- Download the corresponding installation package from [releases](https://github.com/esrrhs/pingtunnel/releases), such as pingtunnel_windows64.zip, and decompress it
- Then run with **administrator** privileges. The commands corresponding to different forwarding functions are as follows.
- If you see a log of ping pong, the connection is normal
- `--key` parameter is **INT** type, only supports numbers between `0-2147483647`

#### Client: Sock5 Proxy

```shell
pingtunnel client -l :1080 -s example.com
```

```shell
NAME:
   pingtunnel client - Run a PingTunnel Client: Sock5 Proxy

USAGE:
   pingtunnel client [command options] [arguments...]

OPTIONS:
   -l address, --socks address   Local address, traffic sent to this port will be forwarded to the server (default: ":1080")
   -s address, --server address  The address of the server, the traffic will be forwarded to this server through the tunnel (default: "exaple.com")
   --key key                     numeric key in [0-2147483647] (default: 0)
   --tcp_bs bytes                Tcp send and receive buffer size (bytes) (default: 1048576)
   --tcp_mw window               The maximum window of tcp (default: 20000)
   --tcp_rst buffer              max process thread's buffer in server (default: 1000)
   --tcp_gz value                tcp will compress data when the packet exceeds this size, 0 means no compression (default: 0)
   --timeout s                   timeout(s) period for the server to initiate a connection to the destination address (default: 60)
```

#### Tunnel: TCP/UDP Port Forward

```shell
NAME:
   pingtunnel tunnel - Run a PingTunnel Tunnel: TCP/UDP Port Forward

USAGE:
   pingtunnel tunnel [command options] [arguments...]

OPTIONS:
   -l address, --listen address  Local address, traffic sent to this port will be forwarded to the server (default: ":5300")
   -t address, --target address  Destination address forwarded by the remote server, traffic will be forwarded to this address (default: "8.8.8.8:53")
   -u, --udp                     forward UDP if defined
   -s address, --server address  The address of the server, the traffic will be forwarded to this server through the tunnel (default: "exaple.com")
   --key key                     numeric key in [0-2147483647] (default: 0)
   --tcp_bs bytes                Tcp send and receive buffer size (bytes) (default: 1048576)
   --tcp_mw window               The maximum window of tcp (default: 20000)
   --tcp_rst buffer              max process thread's buffer in server (default: 1000)
   --tcp_gz value                tcp will compress data when the packet exceeds this size, 0 means no compression (default: 0)
   --timeout s                   timeout(s) period for the server to initiate a connection to the destination address (default: 60)
```

- TCP

  ```shell
  pingtunnel tunnel -l :5300 -t 8.8.8.8:53 -s example.com
  ```

- UDP

  ```shell
  pingtunnel tunnel -l :5300 -t 8.8.8.8:53 -s example.com -u
  ```

### Use Docker

It can also be started directly with docker, which is more convenient. Same parameters as above

- server:

```shell
docker run -d --name ptS --privileged --restart=unless-stopped --network host \
  -e ACTION=server -e ARGS=--key=11223344 \
  pingtunnel
```

- client:

```shell
docker run -d --name ptC --restart=unless-stopped -p 1080:1080 \
  -e ACTION=client -e ARGS='-s=172.17.0.1 --key=11223344' \
  pingtunnel
```

## Thanks for free JetBrains Open Source license

<img src="https://resources.jetbrains.com/storage/products/company/brand/logos/GoLand.png" height="200"/></a>
