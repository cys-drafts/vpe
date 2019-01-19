package tuntap

import (
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"unsafe"
)

type Tap struct {
	Fd int
	Mtu int
	Name string
}

type ifreq struct {
	Name [syscall.IFNAMSIZ]byte
	Flags uint16
}

func ioctl(a1, a2, a3 uintptr) error {
	if _,_,errno := syscall.Syscall(syscall.SYS_IOCTL, a1, a2, a3); errno != 0 {
		return errno
	}
	return nil
}

func run(ifup, name string, mtu int) error {
	cmd := exec.Command(ifup)
	cmd.Env = append(os.Environ(), "IFNAME="+name, "MTU="+strconv.Itoa(mtu))
	if e := cmd.Run(); e != nil {
		return e
	}
	return nil
}

func (tap *Tap) Setup(name string, mtu int, ifup string) error {
	fd, e := syscall.Open("/dev/net/tun", syscall.O_RDWR, 0644)
	if e != nil {
		return e
	}

	var req ifreq

	copy(req.Name[:(syscall.IFNAMSIZ - 1)], name)
	req.Flags = syscall.IFF_TAP | syscall.IFF_NO_PI;
	if e := ioctl(uintptr(fd), uintptr(syscall.TUNSETIFF), uintptr(unsafe.Pointer(&req))); e != nil {
		syscall.Close(fd)
		return e
	}

	if e := run(ifup, name, mtu); e != nil {
		syscall.Close(fd)
		return e
	}

	tap.Fd = fd
	tap.Mtu = mtu
	tap.Name = name

	return nil
}

func (tap *Tap) Read(b []byte) (n int, e error) {
	n, e = syscall.Read(tap.Fd, b)
	return
}

func (tap *Tap) Write(b []byte) (n int, e error) {
	n, e = syscall.Write(tap.Fd, b)
	return
}
