package irc

import (
	"net"
)

type peer interface {
	ID() int
	RemoteAddr() string
	Close() error
	Write([]byte) (int, error)
	ReadAll() ([]byte, error)
}

type Peer struct {
	Id   int
	conn net.Conn
}

func (p *Peer) RemoteAddr() string {
	return p.conn.RemoteAddr().String()
}

func (p *Peer) Write(b []byte) (int, error) {
	return p.conn.Write(b)
}

func (p *Peer) ID() int {
	return p.Id
}

func (p *Peer) ReadAll() ([]byte, error) {
	b := make([]byte, 1024)
	n, err := p.conn.Read(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}

func (p *Peer) Close() error {
	return p.conn.Close()
}
