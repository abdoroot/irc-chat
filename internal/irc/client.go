package irc

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
)

type client struct {
	ServerAddr string
}

func NewClient(option Options) *client {
	addr := ":8080"
	if option.Addr != "" {
		addr = option.Addr
	}
	return &client{ServerAddr: addr}
}

func (c *client) Dial() error {
	conn, err := net.Dial("tcp", c.ServerAddr)
	if err != nil {
		return err
	}

	slog.Info("connected to server ", "addr", c.ServerAddr)

	go func() {
		err := c.SendMessageLoop(conn)
		slog.Error("error sending to server", "err", err)
	}()

	return c.HandleServerResponseLoop(conn)
}

func (c *client) SendMessageLoop(conn net.Conn) error {
	for {
		r := bufio.NewReader(os.Stdin)
		buf := make([]byte, 1024)
		n, err := r.Read(buf)
		if err != nil {
			slog.Error("error reading form stdin", "err", err)
		}
		//Send it to server
		conn.Write(buf[:n])
	}
}

func (c *client) HandleServerResponseLoop(conn net.Conn) (err error) {
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			slog.Error("error reading form server", "err", err)
			if err == io.EOF {
				slog.Info("closing the conn to server ", "addr", c.ServerAddr)
				conn.Close()
				break
			}
		}
		//handle it if needed
		fmt.Println(string(buf[:n]))
	}
	return err
}
