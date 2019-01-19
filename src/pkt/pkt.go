package pkt

import (
	"fmt"
	"encoding/binary"
)

type pktfmt struct {
	Len uint32
}

func Read(b []byte, readfn func([]byte)(int, error)) (int, error) {
	_, e := readfn(b[:4])
	if e != nil {
		return 0, e
	}
	plen := int(binary.BigEndian.Uint32(b[:4]))
	l := 0
	for l < plen {
		n, e := readfn(b[l:plen])
		if e != nil {
			return 0, e
		}
		l += n
	}
	if l != plen {
		fmt.Printf("BUG: len != plen\n")
	}
//	fmt.Printf("pkt read %v bytes\n", l)
	return plen, nil
}

func Write(b []byte, writefn func([]byte)(int, error)) (int, error) {
	plen := len(b)
	b4 := make([]byte, 4)
	binary.BigEndian.PutUint32(b4, uint32(plen))
	_, e := writefn(b4)
	if e != nil {
		return 0, e
	}
	l := 0
	for l < plen {
		n, e := writefn(b[l:plen])
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
