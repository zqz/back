package ws

type Message struct {
	Author string `json:"author"`
	Body   string `json:"body"`
}

type Event struct {
	E string      `json:"e"` // The kind of event.
	P interface{} `json:"p"`
}

func RegisterEvent(id string) *Event {
	return &Event{
		E: "register",
		P: map[string]string{
			"id": id,
		},
	}
}

func (self *Message) String() string {
	return self.Author + " says " + self.Body
}
