package pkt

import (
	"fmt"
	"encoding/binary"
	"net"
	"sock"
)

type pktfmt struct {
	Len uint32
}

func Read(b []byte, s *sock.Sock) (int, net.Addr, error) {
	switch s.Conn.(type) {
		case *net.UDPConn:
			return s.Read(b)
	}

	l := 0
	plen := 4
	for l < plen {
		n, _, e := s.Read(b[:4])
		if e != nil {
			return 0, nil, e
		}
		l +=n
	}
	plen = int(binary.BigEndian.Uint32(b[:4]))
	l = 0
	for l < plen {
		n, _, e := s.Read(b[l:plen])
		if e != nil {
			return 0, nil, e
		}
		l += n
	}
	if l != plen {
		fmt.Printf("BUG: len != plen\n")
	}
//	fmt.Printf("pkt read %v bytes\n", l)
	return plen, nil, nil
}

func Write(b []byte, s *sock.Sock) (int, error) {
	switch s.Conn.(type) {
		case *net.UDPConn:
			return s.Write(b)
	}

	plen := len(b)
	b4 := make([]byte, 4)
	binary.BigEndian.PutUint32(b4, uint32(plen))
	_, e := s.Write(b4)
	if e != nil {
		return 0, e
	}
	l := 0
	for l < plen {
		n, e := s.Write(b[l:plen])
		if e != nil {
			return 0, e
		}
		l += n
	}
	if l != plen {
		fmt.Printf("BUG: len != plen\n")
	}
//	fmt.Printf("pkt write %v bytes\n", l)
	return plen, nil
}
