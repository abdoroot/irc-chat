package main

import (
	"fmt"

	"github.com/abdoroot/irc-chat/internal/irc"
)

func main() {
	s := irc.NewServer(irc.Options{Addr: ":8001"})
	fmt.Println(s.Start())
}
