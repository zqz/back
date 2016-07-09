package p2p

import (
	"fmt"

	"github.com/gorilla/websocket"
)

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
