package ws

import (
	"github.com/zqzca/back/dependencies"
	"golang.org/x/net/websocket"
)

// Chat server.
type Server struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	sendAllCh  chan *Event
	doneCh     chan bool
	errCh      chan error

	*dependencies.Dependencies
}

// Create new chat server.
func NewServer() *Server {
	clients := make(map[string]*Client)
	register := make(chan *Client)
	unregister := make(chan *Client)
	sendAllCh := make(chan *Event)
	doneCh := make(chan bool)
	errCh := make(chan error)

	return &Server{
		clients,
		register,
		unregister,
		sendAllCh,
		doneCh,
		errCh,
		nil,
	}
}

func (s *Server) Register(c *Client) {
	s.register <- c
}

func (s *Server) Unregister(c *Client) {
	s.unregister <- c
}

func (s *Server) SendAll(e *Event) {
	s.sendAllCh <- e
}

func (s *Server) Done() {
	s.doneCh <- true
}

func (s *Server) Err(err error) {
	s.errCh <- err
}

func (s *Server) sendAll(e *Event) {
	for _, c := range s.clients {
		c.Write(e)
	}
}

func (s *Server) WriteClient(cID string, e string, p interface{}) {
	c, ok := s.clients[cID]
	if !ok {
		s.Info("WS: Client ID not found.", "id", cID)
		return
	}
	event := &Event{E: e, P: p}
	s.send(c, event)
}

func (s *Server) Broadcast(e string, p interface{}) {
	for cID, _ := range s.clients {
		s.WriteClient(cID, e, p)
	}
}

func (s *Server) send(c *Client, e *Event) {
	c.Write(e)
}

func (s *Server) Start() {
	s.handleChannels()
}

func (s *Server) Endpoint() websocket.Handler {
	return websocket.Handler(func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				s.errCh <- err
			}
		}()

		client := NewClient(ws, s)
		s.Register(client)
		client.Listen()
	})
}

func (s *Server) handleChannels() {
	s.Info("Websocket Active", "endpoint", "/ws")

	for {
		select {

		// Add new a client
		case c := <-s.register:
			s.clients[c.id] = c
			s.send(c, RegisterEvent(c.id))
			s.Info("WS Added Client", "client", c.id, "total", len(s.clients))

		// del a client
		case c := <-s.unregister:
			s.Info("WS Deleted Client", "client", c.id)
			delete(s.clients, c.id)

		// broadcast message for all clients
		case e := <-s.sendAllCh:
			s.Info("WS Broadcasting to all", "event", e)
			s.sendAll(e)

		case err := <-s.errCh:
			s.Error("WS Error", "error", err.Error())

		case <-s.doneCh:
			s.Info("WS Finishing")
			return
		}
	}
}
