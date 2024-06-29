package irc

import (
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net"
	"slices"
	"sync"
)

var (
	joinChatMsg  = "user %v join chat\n"
	leaveChatMsg = "user %v leave chat\n"
)

type server struct {
	Addr string

	mu    sync.Mutex
	peers []peer
}

type Options struct {
	Addr string
}

func NewServer(option Options) *server {
	addr := ":8080"
	if option.Addr != "" {
		addr = option.Addr
	}
	return &server{
		Addr:  addr,
		peers: []peer{},
	}
}

func (s *server) Start() error {
	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	slog.Info("Irc server running at", "addr", s.Addr)
	return s.HandleConnectionsLoop(l)
}

func (s *server) HandleConnectionsLoop(l net.Listener) error {
	for {
		conn, err := l.Accept()
		if err != nil {
			slog.Error("err accepting conn")
		}

		//add peer to sever
		p := &Peer{conn: conn, Id: RandPeerID()}
		s.mu.Lock()
		s.peers = append(s.peers, p)
		s.mu.Unlock()
		slog.Info("peer added", "peer", p.RemoteAddr(), "id", p.ID())

		//send joining chat
		msg := fmt.Sprintf(joinChatMsg, p.Id)
		s.BroadcastPeerMessages([]byte(msg), p.ID())

		go func() {
			//read from client
			s.ReadClientMessagesLoop(p)
		}()
	}
}

func RandPeerID() int {
	return rand.Intn(500)
}

func (s *server) removePeer(p peer) error {
	if slices.Contains(s.peers, p) {
		for i, v := range s.peers {
			if v.ID() == p.ID() {
				//Delete it
				s.mu.Lock()
				s.peers = append(s.peers[:i], s.peers[i+1:]...)
				s.mu.Unlock()
				slog.Info("peer delete for the server db", "id", p.ID())
				msg := fmt.Sprintf(leaveChatMsg, p.ID())
				s.BroadcastPeerMessages([]byte(msg), p.ID())
				return nil
			}
		}
	}
	return fmt.Errorf("peer not found in server database")
}

func (s *server) ReadClientMessagesLoop(p peer) (err error) {
	slog.Info("reading loop for peer", "id", p.ID())
	for {
		b, err := p.ReadAll()
		if err != nil {
			slog.Error("err accepting conn", "err", err)
			if err == io.EOF {
				slog.Info("closing the conn to peer ", "addr", p.ID())
				err := s.removePeer(p)
				if err != nil {
					slog.Error("error removePeer", "err", err)
				}
				p.Close()
				break
			}
		}
		//handle the client messages
		s.handleClientMessages(b, p)
	}
	return
}

func (s *server) handleClientMessages(b []byte, p peer) {
	//now just prodcast the message
	s.BroadcastPeerMessages(b, p.ID())
}

func (s *server) BroadcastPeerMessages(b []byte, excepted int) error {
	//prodcast the message except the sender
	slog.Info("prodcasting to connected peers")
	counter := 0
	for _, p := range s.peers {
		if p.ID() != excepted {
			//expect the sender
			msg := []byte(fmt.Sprintf("peer %v: ", excepted))
			msg = append(msg, b...)
			_, err := p.Write(msg)
			if err != nil {
				return err
			}
			slog.Info("prodcasting to peer", "id", p.ID())
			counter++
		}
	}

	slog.Info("prodcasting finish", "count", counter)
	return nil
}
