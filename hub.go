package main

type Hub struct {
	broadcast  chan []byte
	clients    map[*Client]struct{}
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]struct{}),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.register:
			hub.clients[client] = struct{}{} // register client
		case client := <-hub.unregister:
			if _, ok := hub.clients[client]; ok { //unregister once
				delete(hub.clients, client) //unregister client
				close(client.send)          //stop sending msg to the closed client
			}
		case msg := <-hub.broadcast:
			for client := range hub.clients {
				client.send <- msg
			}
		}
	}
}
