package transport

import (
	"encoding/json"
	"log"
	"net/http"

	"Harshvardhan2212/websocket_server/internal/realtime"

	"github.com/google/uuid"
)

type WsHandler struct {
	Hub *realtime.Hub
}

func NewWsHnadler(Hub *realtime.Hub) *WsHandler {
	return &WsHandler{
		Hub,
	}
}

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
