package router

import (
	"fmt"
	"log"
	"net/http"

	"Harshvardhan2212/websocket_server/internal/realtime"
	"Harshvardhan2212/websocket_server/internal/transport"

	"github.com/gorilla/mux"
)

type Router struct {
	Mux       *mux.Router
	WsHandler *transport.WsHandler
}

func NewRouter(Hub *realtime.Hub) *Router {
	return &Router{
		Mux:       mux.NewRouter(),
		WsHandler: transport.NewWsHnadler(Hub),
	}
}

func (r *Router) registerRoutes() {
	r.Mux.HandleFunc("/health", transport.HealthCheck)
	r.Mux.HandleFunc("/ws", r.WsHandler.HandleWs)
}

func (r *Router) Run(addr string) {
	r.registerRoutes()
	log.Printf("Listening to port %s", addr)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", addr), r.Mux); err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
