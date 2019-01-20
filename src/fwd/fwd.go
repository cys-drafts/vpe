package fwd

import (
	"fmt"
	"net"
	"encoding/binary"
	"pkt"
	"tuntap"
)

//TODO udp

var fdb = make(map[uint64]*net.TCPConn)

func mac2key(b []byte) uint64 {
	b6 := []byte{b[0],b[1],b[2],b[3],b[4],b[5],0,0}
	return binary.LittleEndian.Uint64(b6)
}

func InsertFDB(b []byte, c *net.TCPConn) {
	key := mac2key(b)
	_, ok := fdb[key]
	if ok {
		fmt.Printf("fdb for key %v updated to %v\n", b, c.RemoteAddr())
	} else {
		fmt.Printf("fdb for key %v inserted to %v\n", b, c.RemoteAddr())
	}
	fdb[key] = c
}

func learnFDB(b []byte, c *net.TCPConn) {
	key := mac2key(b)
	_, ok := fdb[key]
	if !ok {
		fdb[key] = c
		fmt.Printf("fdb for key %v learned to %v\n", b, c.RemoteAddr())
	}
}

func lookupFDB(b []byte) *net.TCPConn {
	key := mac2key(b)
	c, ok := fdb[key]
	if !ok {
		return nil
	}
	return c
}

func invalidateFDB(c *net.TCPConn) {
	for k, v := range fdb {
		if v == c {
			delete(fdb, k)
		}
	}
}

func Fwd2Remote(tap *tuntap.Tap) {
	b := make([]byte, tap.Mtu * 2)

	for {
		plen, e := tap.Read(b)
		if e != nil {
			fmt.Println(e)
			break
		}
//		fmt.Printf("tap read %v bytes\n", plen)
		c := lookupFDB(b[:6])
		if c != nil {
			n, e := pkt.Write(b[:plen], c.Write)
			if e != nil {
				fmt.Println(e)
				c.Close()
				invalidateFDB(c)
			}
			if n != plen {
				fmt.Printf("BUG: sock.write() != plen\n")
			}
		}
	}
}

func Fwd2Local(c *net.TCPConn, tap *tuntap.Tap) {
	b := make([]byte, tap.Mtu * 2)

	for {
		plen, e := pkt.Read(b, c.Read)
		if e != nil {
			fmt.Printf("err 1\n")
			fmt.Println(e)
			c.Close()
			break
		}
		learnFDB(b[6:12], c)
		n, e := tap.Write(b[:plen])
		if e != nil {
			fmt.Printf("err 2\n")
			fmt.Println(e)
			continue
		}
		if n != plen {
			fmt.Printf("BUG: tap.write() != plen\n")
		}
//		fmt.Printf("tap read %v bytes\n", n)
	}
}
