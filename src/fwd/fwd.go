package fwd

import (
	"fmt"
	"encoding/binary"
	"pkt"
	"tuntap"
	"sock"
)

var fdb = make(map[uint64]*sock.Sock)

func mac2key(b []byte) uint64 {
	b6 := []byte{b[0],b[1],b[2],b[3],b[4],b[5],0,0}
	return binary.LittleEndian.Uint64(b6)
}

func InsertFDB(b []byte, s *sock.Sock) {
	key := mac2key(b)
	_, ok := fdb[key]
	if ok {
		fmt.Printf("fdb for key %v updated to %v\n", b, s)
	} else {
		fmt.Printf("fdb for key %v inserted to %v\n", b, s)
	}
	fdb[key] = s
}

func learnFDB(b []byte, s *sock.Sock) {
	key := mac2key(b)
	_, ok := fdb[key]
	if !ok {
		fdb[key] = s
		fmt.Printf("fdb for key %v learned to %v\n", b, s)
	}
}

func lookupFDB(b []byte) *sock.Sock {
	key := mac2key(b)
	s, ok := fdb[key]
	if !ok {
		return nil
	}
	return s
}

func invalidateFDB(s *sock.Sock) {
	for k, v := range fdb {
		if v == s {
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
		s := lookupFDB(b[:6])
		if s != nil {
			n, e := pkt.Write(b[:plen], s)
			if e != nil {
				fmt.Println(e)
				s.Close()
				invalidateFDB(s)
			}
			if n != plen {
				fmt.Printf("BUG: sock.write() != plen\n")
			}
		}
	}
}

func Fwd2Local(c interface{}, tap *tuntap.Tap) {
	b := make([]byte, tap.Mtu * 2)

	s := sock.Sock {
		Conn: c,
		Peer: nil,
	}

	for {
		plen, peer, e := pkt.Read(b, &s)
		if e != nil {
			fmt.Printf("err 1\n")
			fmt.Println(e)
			s.Close()
			break
		}
		learnFDB(b[6:12], &sock.Sock{Conn:c, Peer: peer})
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
