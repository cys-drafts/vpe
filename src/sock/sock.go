package sock

import (
	"net"
	"errors"
)

type Sock struct {
	Conn interface{}
	Peer net.Addr
}

func (s *Sock) Read(b []byte) (int, net.Addr, error) {
	switch s.Conn.(type) {
		case *net.TCPConn:
			ts := s.Conn.(*net.TCPConn)
			l, e := ts.Read(b)
			return l, nil, e
		case *net.UDPConn:
			us := s.Conn.(*net.UDPConn)
			return us.ReadFrom(b)
		default:
			return 0, nil, errors.New("bad protocol")
	}
}
func (s *Sock) Write(b []byte) (int, error) {
	switch s.Conn.(type) {
		case *net.TCPConn:
			ts := s.Conn.(*net.TCPConn)
			return ts.Write(b)
		case *net.UDPConn:
			us := s.Conn.(*net.UDPConn)
			if s.Peer != nil && us.RemoteAddr() == nil {
				r := s.Peer.(net.Addr)
				return us.WriteTo(b, r)
			} else {
				return us.Write(b)
			}
		default:
			return 0, errors.New("bad protocol")
	}
}
func (s *Sock) Close() {
	switch s.Conn.(type) {
		case *net.TCPConn:
			ts := s.Conn.(*net.TCPConn)
			ts.Close()
		case *net.UDPConn:
			us := s.Conn.(*net.UDPConn)
			us.Close()
		default:
			panic("bad protocol")
	}
}

func ConnectTCP(addr string) (*net.TCPConn, error) {
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
func ConnectUDP(addr string) (*net.UDPConn, error) {
	ta, e := net.ResolveUDPAddr("udp", addr)
	if e != nil {
		return nil, e
	}
	c, e := net.DialUDP("udp", nil, ta)
	if e != nil {
		return nil, e
	}
	return c, nil
}
func Connect(proto, addr string) (interface{}, error) {
	if proto == "tcp" {
		return ConnectTCP(addr)
	} else if proto == "udp" {
		return ConnectUDP(addr)
	} else {
		return nil, errors.New("bad protocol")
	}
}

func ListenTCP(addr string) (*net.TCPListener, error) {
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
func ListenUDP(addr string) (*net.UDPConn, error) {
	ta, e := net.ResolveUDPAddr("udp", addr)
	if e != nil {
		return nil, e
	}
	c, e := net.ListenUDP("udp", ta)
	if e != nil {
		return nil, e
	}
	return c, nil
}
func Listen(proto, addr string) (interface{}, error) {
	if proto == "tcp" {
		return ListenTCP(addr)
	} else if proto == "udp" {
		return ListenUDP(addr)
	} else {
		return nil, errors.New("bad protocol")
	}
}

func Accept(c interface{}) (*net.TCPConn, error) {
	tl := c.(*net.TCPListener)
	return tl.AcceptTCP()
}
