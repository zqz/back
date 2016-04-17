package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
)

func (c *wscon) write(mt int, payload []byte) error {
	//c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

func (c *wscon) Unregister() {
	fmt.Println("Unregistering Client")
	for id, s := range api.sessions {
		if s.ws == c {
			delete(api.sessions, id)
			fmt.Println("Deleted Session: ", id)
		}
	}

	api.unregister <- c
	c.ws.Close()
}

func (c *wscon) ReadPump() {
	defer c.Unregister()

	// c.ws.SetReadLimit(maxMessageSize)
	// c.ws.SetReadDeadline(time.Now().Add(pongWait))
	// c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var r io.Reader
		op, r, err := c.ws.NextReader()

		if err != nil {
			fmt.Println("cant go further")
			fmt.Println(err.Error())
			break
		}

		message, err := ioutil.ReadAll(r)

		switch op {
		case websocket.TextMessage:
			var s p2psession

			if err = json.Unmarshal(message, &s); err != nil {
				fmt.Println(err.Error())
				fmt.Println("Failed to decode json: %s", message)
				break
			}

			s.ID = uuid.NewV4().String()
			s.ID = "foo"
			s.ws = c
			s.createdAt = time.Now()

			api.sessions[s.ID] = s

			m := &p2pmsg{
				ID:  s.ID,
				Msg: "created",
			}

			spew.Dump(s)
			spew.Dump(m)
			data, _ := json.Marshal(m)
			c.write(websocket.TextMessage, data)
		case websocket.BinaryMessage:
			fmt.Println("Ignoring binary message")
		default:
			fmt.Println("other")
		}
	}
}

func (c *wscon) WritePump() {
	fmt.Println("writepump start")
	defer func() {
		fmt.Println("writepump end")
		c.ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				fmt.Println("closing")
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type p2p struct {
	register   chan *wscon
	unregister chan *wscon
	sessions   map[string]p2psession
}

var api p2p

func init() {
	api = p2p{
		register:   make(chan *wscon),
		unregister: make(chan *wscon),
		sessions:   make(map[string]p2psession),
	}
}

type wscon struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

type p2psession struct {
	ID   string `json:"id"`
	SDP  string `json:"sdp"`
	Type string `json:"type"`

	createdAt time.Time
	ws        *wscon
}

type p2pmsg struct {
	ID  string `json:"id"`
	Msg string `json:"msg"`
}

func P2PJoinAnswer(c *echo.Context) error {
	id := c.Param("id")

	s, ok := api.sessions[id]

	if ok == false {
		fmt.Println("failed to find session")
		return c.NoContent(http.StatusNotFound)
	}

	var m p2psession

	if err := c.Bind(&m); err != nil {
		fmt.Println(err.Error())
		return err
	}
	data, _ := json.Marshal(m)
	s.ws.write(websocket.TextMessage, data)

	return c.JSON(http.StatusOK, s)
}

func P2PJoin(c *echo.Context) error {
	id := c.Param("id")

	fmt.Println("sessions", len(api.sessions))
	fmt.Println("id")
	spew.Dump(id)

	fmt.Println("searching..")

	s, ok := api.sessions[id]

	if ok == false {
		fmt.Println("failed to find session")
		return c.NoContent(http.StatusNotFound)
	}

	fmt.Println("Found Session with ID:", s.ID)

	return c.JSON(http.StatusOK, s)
}

func P2PWS(c *echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)

	if err != nil {
		log.Println(err)
		return err
	}

	con := &wscon{
		send: make(chan []byte, 256),
		ws:   ws,
	}

	go con.WritePump()
	con.ReadPump()

	api.register <- con

	return nil
}
