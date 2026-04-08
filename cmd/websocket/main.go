package main

import (
	"flag"
	"log"

	"Harshvardhan2212/websocket_server/internal/realtime"
	"Harshvardhan2212/websocket_server/internal/router"

	"github.com/joho/godotenv"
)

var addr = flag.String("port", "8080", "http service address")

func main() {
	flag.Parse()
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	hub := realtime.NewHub()
	go hub.Run()
	r := router.NewRouter(hub)
	r.Run(*addr)
}
