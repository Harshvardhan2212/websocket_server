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
	r.Mux.Use(middleware.Auth)

	adminRoute := r.Mux.PathPrefix("/").Subrouter()
	adminRoute.Use(middleware.RequireRole(realtime.Admin))
	adminRoute.HandleFunc("/create-channel", r.WsHandler.CreateChannel).Methods("POST")
	adminRoute.HandleFunc("/delete-channel/{id}", r.WsHandler.DeleteChannel).Methods("DELETE")

	modarator := r.Mux.PathPrefix("/").Subrouter()
	modarator.Use(middleware.RequireRole(
		realtime.Admin, realtime.Modarator,
	))
	modarator.HandleFunc("/kick-client", r.WsHandler.KickClient).Methods("PUT")
	modarator.HandleFunc("/mute-client", r.WsHandler.MuteClient).Methods("PUT")

	r.Mux.HandleFunc("/join-channel", r.WsHandler.JoinChannel).Methods("POST")
}

func (r *Router) Run(addr string) {
	r.registerRoutes()
	log.Printf("Listening to port %s", addr)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", addr), r.Mux); err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
