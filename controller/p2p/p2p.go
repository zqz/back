package p2p

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

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

// Answer does things
func Answer(c echo.Context) error {
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

// Join does things
func Join(c echo.Context) error {
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

// Signaling does things
func Signaling() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(
			w, r, nil,
		)

		if err != nil {
			log.Println(err)
			return
		}

		con := &wscon{
			send: make(chan []byte, 256),
			ws:   ws,
		}

		go con.WritePump()
		con.ReadPump()

		api.register <- con
	}
}
