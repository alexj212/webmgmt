package webmgmt

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "github.com/gorilla/websocket"
    "github.com/potakhov/loge"
)

const (
    // Time allowed to write a message to the peer.
    writeWait = 10 * time.Second

    // Time allowed to read the next pong message from the peer.
    pongWait = 60 * time.Second

    // Send pings to peer with this period. Must be less than pongWait.
    pingPeriod = (pongWait * 9) / 10

    // Maximum message size allowed from peer.
    maxMessageSize = 512
)

var (
    newline = []byte{'\n'}
    space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

type MessageType int


type ServerMessage struct {
    Authenticated bool   `json:"authenticated"`
    Prompt       string `json:"prompt"`
    PromptColor         string `json:"prompt_color"`
    Response       string `json:"response"`
    Color         string `json:"color"`
}

type ClientMessage struct {
    Payload string `json:"payload"`
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
    app *MgmtApp
    // The websocket connection.
    conn *websocket.Conn
    // Buffered channel of outbound messages.
    send          chan *ServerMessage
    username      string
    authenticated bool
    loginAttempts int
}

func (c *Client) Send(msg *ServerMessage) {
    c.send <- msg
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
    defer func() {
        c.app.hub.unregister <- c
        err := c.conn.Close()
        if err != nil {
            loge.Error("Error calling conn.Close error: %v", err)
        }

    }()
    c.conn.SetReadLimit(maxMessageSize)
    err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
    if err != nil {
        loge.Error("Error calling conn.SetReadDeadline error: %v", err)
    }

    c.conn.SetPongHandler(func(string) error {
        err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
        if err != nil {
            loge.Error("Error calling conn.SetReadDeadline error: %v", err)
        }

        return nil
    })

    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                loge.Error("error: %v", err)
            }
            break
        }
        message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

        msg, err := ConvertBytesToMessage(message)
        if err == nil {
            loge.Info("rx: %v\n", msg.Payload)
            c.handleMessage(msg)
        } else {
            loge.Error("Error converting ws data to json error: %v\n", err)
        }

    }
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        err := c.conn.Close()
        if err != nil {
            loge.Error("Error calling conn.Close error: %v", err)
        }

    }()
    for {
        select {
        case message, ok := <-c.send:
            err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err != nil {
                loge.Error("Error calling conn.SetWriteDeadline error: %v", err)
            }

            if !ok {
                // The hub closed the channel.
                err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                if err != nil {
                    loge.Error("Error calling conn.WriteMessage error: %v", err)
                }

                return
            }

            err = c.conn.WriteJSON(message)
            if err != nil {
                return
            }
            // Add queued chat messages to the current websocket message.
            n := len(c.send)
            for i := 0; i < n; i++ {
                err := c.conn.WriteJSON(message)
                if err != nil {
                    return
                }
            }

        case <-ticker.C:
            err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err != nil {
                loge.Error("Error calling conn.SetWriteDeadline error: %v", err)
            }

            if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}

func (c *Client) handleMessage(message *ClientMessage) {

    if ! c.authenticated {

        if c.username == "" {
            loge.Info("handleMessage  - prompt password")
            c.authenticated = true
            c.username = message.Payload
            loginPasswordMessage := &ServerMessage{Authenticated:false, Prompt:"Enter Password: ", PromptColor:"green"}
            c.Send(loginPasswordMessage)
            return
        }


        if c.username == "alex" && message.Payload == "bambam" {
            loge.Info("handleMessage  - valid password")
            c.authenticated = true
            cmdPromptMessage := &ServerMessage{Authenticated:true, Prompt:"$ ", PromptColor:"green"}
            c.Send(cmdPromptMessage)
            return
        }

        c.authenticated = false
        c.loginAttempts++

        if c.loginAttempts<3 {
            loge.Info("handleMessage  - invalid password attempts < 3")
            loginUsernameMessage := &ServerMessage{Authenticated:false, Prompt:"Enter Username: ", PromptColor:"green"}
            loginUsernameMessage.Response = "Invalid Password"
            loginUsernameMessage.Color = "red"
            c.Send(loginUsernameMessage)
            return
        }

        loge.Info("handleMessage  - invalid password attempts >= 3")

        loginUsernameMessage := &ServerMessage{Authenticated:false, Prompt:"", PromptColor:""}
        loginUsernameMessage.Response = "Invalid Password - disconnecting client"
        loginUsernameMessage.Color = "red"
        c.Send(loginUsernameMessage)
        time.Sleep( 500 * time.Millisecond)
        _ = c.conn.Close()

        return
    }

    loge.Info("handleMessage  - authenticated user message.Payload: "+message.Payload)
    responseMesg := &ServerMessage{Authenticated:true, Prompt:"$ ", PromptColor:"green"}
    responseMesg.Response = fmt.Sprintf("echo: %v", message.Payload)
    responseMesg.Color = "green"
    c.Send(responseMesg)
}

// func (c *Client) handleMessage(message *Message) {
//
//     in := &api.AdminMessage{}
//     in.Command = message.Payload
//
//     data := fmt.Sprintf("Executing Cmd: %v\n", in.Command)
//     response := &Message{Payload: data, MessageType: output}
//     c.app.hub.Broadcast(response)
//
//     // loge.Info("Got command: %v\n", in.Command)
//     ctx := context.Background()
//     resp, err := c.app.adminServer.Execute(ctx, in)
//     if err != nil {
//         data := fmt.Sprintf("Got Error: %v\n", err)
//         response := &Message{Payload: data, MessageType: output}
//         c.app.hub.Broadcast(response)
//     } else {
//         if resp.Status == api.AdminResponse_failed {
//             data := fmt.Sprintf("Got AdminResponse_failed Error: %v\n", resp.Error)
//             response := &Message{Payload: data, MessageType: output}
//             c.app.hub.Broadcast(response)
//         } else {
//             data := fmt.Sprintf("Got AdminResponse_success Response\n%v\n", resp.Response)
//             response := &Message{Payload: data, MessageType: output}
//             c.app.hub.Broadcast(response)
//         }
//
//     }
// }

// serveWs handles websocket requests from the peer.
func serveWs(app *MgmtApp, w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        loge.Error("serveWs error: %v\n", err)
        return
    }
    client := &Client{app: app, conn: conn, send: make(chan *ServerMessage, 256)}
    client.app.hub.register <- client

    // Allow collection of memory referenced by the caller by doing all work in
    // new goroutines.
    go client.writePump()
    go client.readPump()


    loginUsernameMessage := &ServerMessage{Authenticated:false, Prompt:"Enter Username: ", PromptColor:"green"}
    loginUsernameMessage.Response = "Welcome to the server"
    loginUsernameMessage.Color = "red"
    client.Send(loginUsernameMessage)

}

func ConvertBytesToMessage(payload []byte) (*ClientMessage, error) {
    msg := &ClientMessage{}
    err := json.Unmarshal(payload, &msg)
    return msg, err
}
