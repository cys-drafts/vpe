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
		conn, e := sock.Listen(conf.Server)
		if e != nil {
			fmt.Println(e)
			return
		}
		fmt.Printf("setup listener\n")
		for {
			c, e := conn.AcceptTCP()
			if e != nil {
				fmt.Println(e)
				continue
			}
			//fmt.Printf("incoming connection %v\n", c)
			go fwd.Fwd2Local(c, tap)
		}
	} else {
		conn, e := sock.Connect(conf.Server)
		if e != nil {
			fmt.Println(e)
			return
		}
		//fmt.Printf("connected to %v\n", conf.Server)
		fwd.InsertFDB([]byte{0xff,0xff,0xff,0xff,0xff,0xff}, conn)
		fwd.Fwd2Local(conn, tap)
	}
}
