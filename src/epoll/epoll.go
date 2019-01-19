package epoll

import (
	"syscall"
)

type Epoll struct {
	Epfd int
}

func Create() (*Epoll, error) {
	fd, e := syscall.EpollCreate(1)
	if e != nil {
		return nil, e
	}
	ep := Epoll{ Epfd : fd }
	return &ep, nil
}

func (ep *Epoll) AddIn(fd int) error {
	event := syscall.EpollEvent{
		Events : syscall.EPOLLIN,
		Fd : int32(fd),
	}
	return syscall.EpollCtl(ep.Epfd, syscall.EPOLL_CTL_ADD, fd, &event)
}

func (ep *Epoll) Poll(events []EpollEvent, msec int) (int, error) {
	return syscall.EpollWait(ep.Epfd, events, msec)
}
