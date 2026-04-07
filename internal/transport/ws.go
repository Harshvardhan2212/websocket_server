package transport

import (
	"encoding/json"
	"log"
	"net/http"

	"Harshvardhan2212/websocket_server/internal/realtime"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type WsHandler struct {
	Hub *realtime.Hub
}

func NewWsHnadler(Hub *realtime.Hub) *WsHandler {
	return &WsHandler{
		Hub,
	}
}

func (ws *WsHandler) CreateChannel(w http.ResponseWriter, r *http.Request) {
	channID := ws.Hub.CreateChannel()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"channel_id": channID.String(),
		"message":    "channel created successfully",
	})
}

func (ws *WsHandler) DeleteChannel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	channID, err := uuid.Parse(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid UUID provided",
		})
		return
	}

	ok := ws.Hub.DeleteChannel(channID)

	if ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"message": "channel deleted successfully",
		})
	}
}

func (ws *WsHandler) JoinChannel(w http.ResponseWriter, r *http.Request) {}

func (ws *WsHandler) KickClient(w http.ResponseWriter, r *http.Request) {}

func (ws *WsHandler) MuteClient(w http.ResponseWriter, r *http.Request) {}

func (ws *WsHandler) HandleWs(w http.ResponseWriter, r *http.Request) {
	var channID uuid.UUID

	id := r.URL.Query().Get("id")

	if id == "" {
		channID = uuid.New()
	} else {
		parsedID, err := uuid.Parse(id)
		if err != nil {
			log.Println("Invalid UUID:", err)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]string{
				"message": "Invalid UUID provided",
			})
			return // ✅ IMPORTANT
		}
		channID = parsedID
	}

	realtime.ServeWs(ws.Hub, w, r, channID)
}
