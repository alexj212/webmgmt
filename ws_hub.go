package webmgmt

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*WSClient]bool

	// Inbound messages from the clients.
	broadcast chan ServerMessage

	// Register requests from the clients.
	register chan *WSClient

	// Unregister requests from clients.
	unregister chan *WSClient
}

// newHub will create the struct for handling clients.
func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan ServerMessage),
		register:   make(chan *WSClient),
		unregister: make(chan *WSClient),
		clients:    make(map[*WSClient]bool),
	}
}

// run handler for processing new clients and client disconnects and broadcast ServerMessages.
func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				client.connected = false
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					client.connected = false
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// Broadcast sill send a ServerMessage to all users on the System.
func (h *Hub) Broadcast(msg ServerMessage) {
	h.broadcast <- msg
}
