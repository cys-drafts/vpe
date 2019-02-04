package main

import (
	"fmt"
	"config"
	"fwd"
	"tuntap"
	"sock"
)

func main() {
	conf, e := config.Parse()
	if e != nil {
		fmt.Println(e)
		return
	}
	fmt.Println(conf)

	tap := &tuntap.Tap{}
	e = tap.Setup(conf.Ifname, conf.Mtu, conf.Ifscript)
	if e != nil {
		fmt.Println(e)
		return
	}
	fmt.Printf("setup %s\n", conf.Ifname)

	go fwd.Fwd2Remote(tap)

	if conf.Role == config.ROLE_SERVER {
		conn, e := sock.Listen(conf.Proto, conf.Server)
		if e != nil {
			fmt.Println(e)
			return
		}
		fmt.Printf("setup listener\n")
		if conf.Proto == "tcp" {
			for {
				c, e := sock.Accept(conn)
					if e != nil {
						fmt.Println(e)
							continue
					}
				//fmt.Printf("incoming connection %v\n", c)
				go fwd.Fwd2Local(c, tap)
			}
		} else if conf.Proto == "udp" {
			fwd.Fwd2Local(conn, tap)
		} else {
			fmt.Printf("bad protocol")
		}
	} else {
		conn, e := sock.Connect(conf.Proto, conf.Server)
		if e != nil {
			fmt.Println(e)
			return
		}
		fmt.Printf("connected to %v\n", conf.Server)
		fwd.InsertFDB([]byte{0xff,0xff,0xff,0xff,0xff,0xff}, &sock.Sock{Conn:conn, Peer:nil,})
		fwd.Fwd2Local(conn, tap)
	}
}
