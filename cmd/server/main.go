package main

import (
	"flag"

	"Harshvardhan2212/websocket_server/internal/realtime"
	"Harshvardhan2212/websocket_server/internal/router"
)

var addr = flag.String("port", ":8080", "http service address")

func main() {
	flag.Parse()

	hub := realtime.NewHub()
	go hub.Run()
	r := router.NewRouter(hub)
	r.Run(*addr)
}
