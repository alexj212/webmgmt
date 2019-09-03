package webmgmt

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
    // Registered clients.
    Clients map[*Client]bool

    // Inbound messages from the clients.
    broadcast chan *ServerMessage

    // Register requests from the clients.
    register chan *Client

    // Unregister requests from clients.
    unregister chan *Client
}

func newHub() *Hub {
    return &Hub{
        broadcast:  make(chan *ServerMessage),
        register:   make(chan *Client),
        unregister: make(chan *Client),
        Clients:    make(map[*Client]bool),
    }
}

func (h *Hub) run() {
    for {
        select {
        case client := <-h.register:
            h.Clients[client] = true
        case client := <-h.unregister:
            if _, ok := h.Clients[client]; ok {
                client.connected = false
                delete(h.Clients, client)
                close(client.send)
            }
        case message := <-h.broadcast:
            for client := range h.Clients {
                select {
                case client.send <- message:
                default:
                    client.connected = false
                    close(client.send)
                    delete(h.Clients, client)
                }
            }
        }
    }
}

func (h *Hub) Broadcast(msg *ServerMessage) {
    h.broadcast <- msg
}
