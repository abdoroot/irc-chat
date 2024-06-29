package main

import (
	"fmt"

	"github.com/abdoroot/irc-chat/internal/irc"
)

func main() {
	// client
	c := irc.NewClient(irc.Options{Addr: "127.0.0.1:8001"})
	fmt.Println(c.Dial())
}
