package eventHandler

import (
	"bytes"
	"context"
	"encoding/json"
	"git.raad.cloud/cloud/hermes/pkg/auth"
	"git.raad.cloud/cloud/hermes/pkg/newMessage"
	"git.raad.cloud/cloud/hermes/pkg/read"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 100000
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
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.Unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		err = c.handleEvents(message)
		if err != nil {
			logrus.Errorf("error in handling event: %v", err)
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

func (c *Client) handleEvents(message []byte) error {
	logrus.Info("inja")
	event := map[string]interface{}{}

	logrus.Println(string(message))
	err := json.Unmarshal(message, &event)
	if err != nil {
		return errors.Wrap(err, "erorr while trying to unmarshall message from ws")
	}
	token, ok := event["token"]
	if !ok {
		return errors.New("token not found in event")
	}
	ident, err := auth.GetAuthentication(token.(string), "")
	if err != nil {
		return errors.Wrap(err, "could not authorize token")
	}
	logrus.Info("Token is OK")
	BaseHub.ClientsMap[ident.ID] = c
	eventType, ok := event["type"]
	if !ok {
		return errors.Wrap(err, "no event type in event found")
	}
	switch eventType.(string) {
	case "DLV":
		jp := &JoinPayload{}
		err := mapstructure.Decode(event, jp)
		if err != nil {
			return errors.Wrap(err, "error while decoding event into NewMessage")
		}
		Handle(context.Background(), jp)
		//if err != nil {
		//	logrus.Errorf("Error in Join event : %v", err)
		//}

	case "NEW":
		logrus.Info("Event is New Message")
		nm := &newMessage.NewMessage{}
		err := mapstructure.Decode(event, nm)
		if err != nil {
			return errors.Wrap(err, "error while decoding event into NewMessage")
		}
		err = newMessage.Handle(nm)
		if err != nil {
			logrus.Errorf("Error in NewMessage Event : %v", err)
		}

	case "READ":
		logrus.Info("Event is read")
		rs := &read.ReadSignal{}
		err := mapstructure.Decode(event, rs)
		if err != nil {
			return errors.Wrap(err, "error while decoding event into ReadSignal")
		}
		err = read.Handle(rs)
		if err != nil {
			logrus.Errorf("Error in handling read signal")
		}
	case "JOIN":
		logrus.Println(" event is JOin")
		rs := &JoinPayload{}
		err := mapstructure.Decode(event, rs)
		if err != nil {
			return errors.Wrap(err, "error while decoding event into ReadSignal")
		}
		rs.UserID = ident.ID
		Handle(context.Background(), rs)
	}
	return nil
}
