package sock

import (
	"net"
)

//TODO udp

func Connect(addr string) (*net.TCPConn, error) {
	ta, e := net.ResolveTCPAddr("tcp", addr)
	if e != nil {
		return nil, e
	}
	c, e := net.DialTCP("tcp", nil, ta)
	if e != nil {
		return nil, e
	}
	return c, nil
}

func Listen(addr string) (*net.TCPListener, error) {
	ta, e := net.ResolveTCPAddr("tcp", addr)
	if e != nil {
		return nil, e
	}
	c, e := net.ListenTCP("tcp", ta)
	if e != nil {
		return nil, e
	}
	return c, nil
}

func Accept(tl *net.TCPListener) (*net.TCPConn, error) {
	return tl.AcceptTCP()
}
