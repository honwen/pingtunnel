package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/esrrhs/pingtunnel"
	"github.com/urfave/cli"
)

var (
	app *cli.App

	working = false

	VersionString = "MISSING build version [git hash]"
)

func init() {
	// -------------------------------- clientFlags --------------------------------
	clientServerFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "s, server",
			Value: "exaple.com",
			Usage: "The `address` of the server, the traffic will be forwarded to this server through the tunnel",
		},
		cli.UintFlag{
			Name:  "key",
			Usage: "numeric `key` in [0-2147483647]",
		},
	}
	clientTCPFlags := []cli.Flag{
		cli.UintFlag{
			Name:  "tcp_bs",
			Value: 1 * 1024 * 1024,
			Usage: "Tcp send and receive buffer size (`bytes`)",
		},
		cli.UintFlag{
			Name:  "tcp_mw",
			Value: 20000,
			Usage: "The maximum `window` of tcp",
		},
		cli.UintFlag{
			Name:  "tcp_rst",
			Value: 1000,
			Usage: "max process thread's `buffer` in server",
		},
		cli.UintFlag{
			Name:  "tcp_gz",
			Usage: "tcp will compress data when the packet exceeds this size, 0 means no compression",
		},
		cli.UintFlag{
			Name:  "timeout",
			Value: 60,
			Usage: "timeout(`s`) period for the server to initiate a connection to the destination address",
		},
	}

	socksFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "l, socks",
			Value: ":1080",
			Usage: "Local `address`, traffic sent to this port will be forwarded to the server",
		}}
	socksFlags = append(socksFlags, clientServerFlags...)
	socksFlags = append(socksFlags, clientTCPFlags...)

	tunnelFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "l, listen",
			Value: ":5300",
			Usage: "Local `address`, traffic sent to this port will be forwarded to the server",
		},
		cli.StringFlag{
			Name:  "t, target",
			Value: "8.8.8.8:53",
			Usage: "Destination `address` forwarded by the remote server, traffic will be forwarded to this address",
		},
		cli.BoolFlag{
			Name:  "u, udp",
			Usage: "forward UDP if defined",
		},
	}
	tunnelFlags = append(tunnelFlags, clientServerFlags...)
	tunnelFlags = append(tunnelFlags, clientTCPFlags...)
	// -------------------------------- clientFlags --------------------------------

	app = cli.NewApp()
	app.Name = "pingtunnel"
	app.Usage = "A tool that send TCP/UDP traffic over ICMP"
	app.Version = fmt.Sprintf("Git:[%s] (%s)", strings.ToUpper(VersionString), runtime.Version())
	app.Commands = []cli.Command{
		{
			Name:  "server",
			Usage: "Run a PingTunnel Server",
			Flags: []cli.Flag{
				cli.UintFlag{
					Name:  "key",
					Usage: "numeric `key` in [0-2147483647]",
				},
				cli.UintFlag{
					Name:  "maxconn",
					Usage: "max num of `connection`s, 0 means no limit",
				},
				cli.UintFlag{
					Name:  "maxthread",
					Value: 100,
					Usage: "max process `thread` in server",
				},
				cli.UintFlag{
					Name:  "maxbuffer",
					Value: 1000,
					Usage: "max process thread's `buffer` in server",
				},
				cli.UintFlag{
					Name:  "timeout",
					Value: 1000,
					Usage: "timeout(`ms`) period for the server to initiate a connection to the destination address",
				},
			},
			Action: func(c *cli.Context) error {
				worker, err := pingtunnel.NewServer(
					c.Int(`key`), c.Int(`maxconn`), c.Int(`maxthread`), c.Int(`maxbuffer`), c.Int(`timeout`),
				)
				if err != nil {
					return err
				}
				working = true
				return worker.Run()
			},
		},
		{
			Name:  "client",
			Usage: "Run a PingTunnel Client: Sock5 Proxy",
			Flags: socksFlags,
			Action: func(c *cli.Context) error {
				filter := func(addr string) bool {
					return true
				}

				if c.Int(`tcp_mw`)*10 > pingtunnel.FRAME_MAX_ID {
					return fmt.Errorf("set tcp win to big, max = %d", pingtunnel.FRAME_MAX_ID/10)
				}

				worker, err := pingtunnel.NewClient(
					c.String(`socks`),
					c.String(`server`),
					"",
					c.Int(`timeout`),
					c.Int(`key`),
					int(1), c.Int(`tcp_bs`), c.Int(`tcp_mw`), c.Int(`tcp_rst`), c.Int(`tcp_gz`), int(0),
					int(1), int(0), &filter)
				if err != nil {
					return err
				}

				working = true
				return worker.Run()
			},
		},
		{
			Name:  "tunnel",
			Usage: "Run a PingTunnel Tunnel: TCP/UDP Port Forward",
			Flags: tunnelFlags,
			Action: func(c *cli.Context) (err error) {
				filter := func(addr string) bool {
					return true
				}

				if c.Int(`tcp_mw`)*10 > pingtunnel.FRAME_MAX_ID {
					return fmt.Errorf("set tcp win to big, max = %d", pingtunnel.FRAME_MAX_ID/10)
				}

				var worker *pingtunnel.Client
				if c.BoolT(`udp`) {
					worker, err = pingtunnel.NewClient(
						c.String(`listen`),
						c.String(`server`),
						c.String(`target`),
						c.Int(`timeout`),
						c.Int(`key`),
						int(0), int(0), int(0), int(0), int(0), int(0),
						int(0), int(0), &filter)
				} else {
					worker, err = pingtunnel.NewClient(
						c.String(`listen`),
						c.String(`server`),
						c.String(`target`),
						c.Int(`timeout`),
						c.Int(`key`),
						int(1), c.Int(`tcp_bs`), c.Int(`tcp_mw`), c.Int(`tcp_rst`), c.Int(`tcp_gz`), int(0),
						int(0), int(0), &filter)
				}
				if err != nil {
					return err
				}

				working = true
				return worker.Run()
			},
		},
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		log.Printf("ERROR: %s", err.Error())
		os.Exit(1)
	}
	for working {
		time.Sleep(time.Hour)
	}
}
