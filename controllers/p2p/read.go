package p2p

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
)

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
				fmt.Println("Failed to decode json: ", message)
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
