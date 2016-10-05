package ws

import (
	"fmt"
	"io"
	"log"

	"github.com/davecgh/go-spew/spew"
	uuid "github.com/satori/go.uuid"

	"golang.org/x/net/websocket"
)

const channelBufSize = 100

// Chat client.
type Client struct {
	id     string
	ws     *websocket.Conn
	server *Server
	ch     chan *Event
	doneCh chan bool
}

// Create new chat client.
func NewClient(ws *websocket.Conn, server *Server) *Client {
	if ws == nil {
		panic("ws cannot be nil")
	}

	if server == nil {
		panic("server cannot be nil")
	}

	id := uuid.NewV4().String()
	ch := make(chan *Event, channelBufSize)
	doneCh := make(chan bool)

	return &Client{id, ws, server, ch, doneCh}
}

func (c *Client) Conn() *websocket.Conn {
	return c.ws
}

func (c *Client) Write(e *Event) {
	select {
	case c.ch <- e:
	default:
		c.server.Del(c)
		err := fmt.Errorf("client %s is disconnected", c.id)
		c.server.Err(err)
	}
}

func (c *Client) Done() {
	c.doneCh <- true
}

// Listen Write and Read request via chanel
func (c *Client) Listen() {
	go c.listenWrite()
	c.listenRead()
}

// Listen write request via chanel
func (c *Client) listenWrite() {
	log.Println("Listening write to client")
	for {
		select {

		// send message to the client
		case e := <-c.ch:
			c.server.Info("WS Client Send", "e", e)
			websocket.JSON.Send(c.ws, e)

		// receive done request
		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true // for listenRead method
			return
		}
	}
}

// Listen read request via chanel
func (c *Client) listenRead() {
	for {
		select {

		// receive done request
		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true // for listenWrite method
			return

		// read data from websocket connection
		default:
			var e Event
			err := websocket.JSON.Receive(c.ws, &e)
			if err == io.EOF {
				c.doneCh <- true
			} else if err != nil {
				c.server.Err(err)
			} else {
				spew.Dump(e)
			}
		}
	}
}
