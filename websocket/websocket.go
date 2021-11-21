package websocket

import (
	"encoding/json"
	"github.com/qownnotes/qc/config"
	"github.com/qownnotes/qc/entity"
	"log"
	"strconv"

	// "net/http"
	"net/url"
	// "strings"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	Token string `json:"token"`
	Type  string `json:"type"`
}

type ResultMessage struct {
	Type            string `json:"type"`
	CommandSnippets []entity.SnippetInfo `json:"data"`
	// Data 			string `json:"data"`
}

// type CommandSnippet struct {
// 	Command      string `json:"command"`
// 	Description  string `json:"description"`
// 	Tags  		 []string `json:"tags"`
// }

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// 20MB maximum message size allowed from peer.
	maxMessageSize = 20971520
)

func FetchSnippetsData() []entity.SnippetInfo {
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:" + strconv.Itoa(config.Conf.QOwnNotes.WebSocketPort)}
	log.Printf("Connecting to QOwnNotes on %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Couldn't connect to QOwnNotes websocket, did you enable the socket server? Error: ", err)
	}
	defer c.Close()

	message := Message{
		Token: config.Conf.QOwnNotes.Token,
		Type:  "getCommandSnippets",
	}

	m, err := json.Marshal(message)

	err = c.WriteMessage(websocket.TextMessage, m)
	if err != nil {
		log.Fatal("Couldn't send command to QOwnNotes: ", err)
	}

	_, msg, err := c.ReadMessage()
	if err != nil {
		log.Fatalf("Couldn't read message from QOwnNotes: %v", err)
	}

	// log.Printf("msg: %v", msg)

	var resultMessage ResultMessage
	err = json.Unmarshal(msg, &resultMessage)
	if err != nil {
		log.Fatalf("Couldn't interpret message from QOwnNotes: %v", err)
	}

	switch resultMessage.Type {
	case "tokenQuery":
		log.Fatal("Please execute \"qc configure\" and configure your token from QOwnNotes!")
	case "commandSnippets":
		log.Printf("CommandSnippets: %v", resultMessage.CommandSnippets)
		return resultMessage.CommandSnippets
	default:
		log.Fatal("Did not understand response from QOwnNotes!")
	}

	// TODO: how to handle a timeout?
	return []entity.SnippetInfo{}
}

//
// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// 	CheckOrigin: func(r *http.Request) bool {
// 		origin := r.Header["Origin"]
//
// 		if len(origin) == 0 {
// 			return true
// 		}
//
// 		u, err := url.Parse(origin[0])
// 		if err != nil {
// 			return false
// 		}
//
// 		// Allow development from other localhost connections
// 		if strings.Contains(u.Host, "localhost:") {
// 			return true
// 		}
//
// 		return strings.ToLower(u.Host) == strings.ToLower(r.Host)
// 	},
// }
//
// // connection is a middleman between the websocket connection and the hub.
// type connection struct {
// 	// The websocket connection.
// 	ws *websocket.Conn
//
// 	// Buffered channel of outbound messages.
// 	send chan []byte
// }
//
// // readPump pumps messages from the websocket connection to the hub.
// func (s subscription) readPump() {
// 	c := s.conn
// 	defer func() {
// 		h.unregister <- s
// 		c.ws.Close()
// 	}()
// 	c.ws.SetReadLimit(maxMessageSize)
// 	c.ws.SetReadDeadline(time.Now().Add(pongWait))
// 	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
// 	for {
// 		_, msg, err := c.ws.ReadMessage()
// 		if err != nil {
// 			// log.Printf("other error: %v", err)
// 			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
// 				log.Printf("error: %v", err)
// 			}
// 			break
// 		}
// 		m := message{msg, s.room, c}
//
// 		// log.Printf("Got message: %#v\n", m)
// 		log.Printf("Got message in room %v", s.room)
//
// 		h.broadcast <- m
// 	}
// }
//
// // write writes a message with the given message type and payload.
// func (c *connection) write(mt int, payload []byte) error {
// 	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
// 	return c.ws.WriteMessage(mt, payload)
// }
//
// // writePump pumps messages from the hub to the websocket connection.
// func (s *subscription) writePump() {
// 	c := s.conn
// 	ticker := time.NewTicker(pingPeriod)
// 	defer func() {
// 		ticker.Stop()
// 		c.ws.Close()
// 	}()
// 	for {
// 		select {
// 		case message, ok := <-c.send:
// 			if !ok {
// 				c.write(websocket.CloseMessage, []byte{})
// 				return
// 			}
// 			if err := c.write(websocket.TextMessage, message); err != nil {
// 				return
// 			}
// 		case <-ticker.C:
// 			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
// 				return
// 			}
// 		}
// 	}
// }
//
// // serveWs handles websocket requests from the peer.
// func serveWs(w http.ResponseWriter, r *http.Request, room string) {
// 	ws, err := upgrader.Upgrade(w, r, nil)
//
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	c := &connection{send: make(chan []byte, 256), ws: ws}
// 	s := subscription{c, room}
// 	h.register <- s
// 	go s.writePump()
// 	go s.readPump()
// }
//
