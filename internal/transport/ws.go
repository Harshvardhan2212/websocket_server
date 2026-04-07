package transport

import (
	"encoding/json"
	"net/http"

	"Harshvardhan2212/websocket_server/internal/middleware"
	"Harshvardhan2212/websocket_server/internal/realtime"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type MuteRequest struct {
	ID string `json:"id"`
}

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
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
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

func (ws *WsHandler) KickClient(w http.ResponseWriter, r *http.Request) {
	var req MuteRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	channID, err := uuid.Parse(req.ID)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	clientVal := r.Context().Value(middleware.ClientID)

	clientID, err := uuid.Parse(clientVal.(string))
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}
	if ws.Hub.KickClient(channID, clientID) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"message": "User kicked from channel",
		})
	}
}

func (ws *WsHandler) MuteClient(w http.ResponseWriter, r *http.Request) {
	var req MuteRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	channID, err := uuid.Parse(req.ID)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	clientVal := r.Context().Value(middleware.ClientID)

	clientID, err := uuid.Parse(clientVal.(string))
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}
	if ws.Hub.MuteClient(channID, clientID) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"message": "User Muted from channel",
		})
	}
}

func (ws *WsHandler) JoinChannel(w http.ResponseWriter, r *http.Request) {
	var req MuteRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	channID, err := uuid.Parse(req.ID)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	clientVal := r.Context().Value(middleware.ClientID)

	clientID, err := uuid.Parse(clientVal.(string))
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	role := r.Context().Value(middleware.Role)

	realtime.ServeWs(ws.Hub, w, r, channID, clientID, realtime.RoleName(role.(string)))
}
