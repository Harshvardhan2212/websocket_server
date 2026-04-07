package router

import (
	"fmt"
	"log"
	"net/http"

	"Harshvardhan2212/websocket_server/internal/middleware"
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

	protected := r.Mux.PathPrefix("/").Subrouter()
	protected.Use(middleware.Auth)
	protected.HandleFunc("/create-channel", r.WsHandler.CreateChannel).Methods("POST")
	protected.HandleFunc("/delete-channel/{id}", r.WsHandler.DeleteChannel).Methods("DELETE")
	protected.HandleFunc("/join-channel", r.WsHandler.JoinChannel).Methods("POST")
	protected.HandleFunc("/kick-client", r.WsHandler.KickClient).Methods("PUT")
	protected.HandleFunc("/mute-client", r.WsHandler.MuteClient).Methods("PUT")
	protected.HandleFunc("/ws", r.WsHandler.HandleWs)
}

func (r *Router) Run(addr string) {
	r.registerRoutes()
	log.Printf("Listening to port %s", addr)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", addr), r.Mux); err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
