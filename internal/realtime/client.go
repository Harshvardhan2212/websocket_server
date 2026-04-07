package realtime

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	ID   uuid.UUID
	hub  *Hub
	conn *websocket.Conn
	send chan *Message
	Role *Role
}

type Register struct {
	ChannelID uuid.UUID
	Client    *Client
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (r *Register) readPump() {
	log.Println("read is running")
	defer func() {
		r.Client.hub.unregister <- r.Client
		r.Client.conn.Close()
	}()
	r.Client.conn.SetReadLimit(maxMessageSize)
	r.Client.conn.SetReadDeadline(time.Now().Add(pongWait))
	r.Client.conn.SetPongHandler(func(string) error { r.Client.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := r.Client.conn.ReadMessage()
		log.Printf("message from connection %v\n", string(message))
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			log.Printf("error from message reading %v\n", err)
			break
		}
		var msg Message

		err = json.Unmarshal(message, &msg)
		if err != nil {
			r.Client.send <- &Message{
				ChannelID: r.ChannelID,
				Payload:   "invalid Message structur",
			}
			log.Println("invalid message:", err)
			continue
		}

		log.Printf("Message is from read : %+v\n", msg)
		log.Println("message reach read and send to broadcast channel")
		r.Client.hub.broadcast <- &msg
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (r *Register) writePump() {
	log.Println("write is running")
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		r.Client.conn.Close()
	}()
	for {
		select {
		case message, ok := <-r.Client.send:
			r.Client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				r.Client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := r.Client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				r.Client.send <- &Message{
					ChannelID: r.ChannelID,
					Payload:   "json parsing issue",
				}
				log.Println("json parsing issue", err)
				continue
			}

			w.Write(data)

			// Add queued chat messages to the current websocket message.
			n := len(r.Client.send)
			for range n {
				w.Write(newline)
				data, err := json.Marshal(<-r.Client.send)
				if err != nil {
					r.Client.send <- &Message{
						ChannelID: r.ChannelID,
						Payload:   "json parsing issue",
					}
					log.Println("json parsing issue", err)
					continue
				}
				w.Write(data)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			r.Client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := r.Client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

/*
* NOTE: this is a entry point which upgrader from http to ws
 */
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request, channID uuid.UUID) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan *Message, 256),
	}
	reg := &Register{
		ChannelID: channID,
		Client:    client,
	}

	client.hub.register <- reg
	log.Printf("client connected to channel id: %v\n", channID)

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go reg.writePump()
	go reg.readPump()
}
