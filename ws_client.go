package webmgmt

import (
    "bytes"
    "encoding/json"
    "net/http"
    "strings"
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

type Client interface {
    IsAuthenticated() bool
    IsConnected() bool
    Username() string
    Send(msg ServerMessage)
    History() []string
    HttpReq() *http.Request
    Misc() map[string]interface{}
}

// Client is a middleman between the websocket connection and the hub.
type WSClient struct {
    app           *MgmtApp
    conn          *websocket.Conn    // The websocket connection.
    send          chan ServerMessage // Buffered channel of outbound messages.
    username      string
    authenticated bool
    loginAttempts int
    connected     bool
    misc          map[string]interface{}
    httpReq       *http.Request
    history       []string
}

func (c *WSClient) Send(msg ServerMessage) {
    if c.connected {
        c.send <- msg
    }
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *WSClient) readPump() {
    defer func() {

        c.app.hub.unregister <- c
        c.conn.Close()
        c.connected = false
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
        c.connected = true
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
func (c *WSClient) writePump() {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        c.conn.Close()
        c.connected = false

        c.app.Config.UnregisterUser(c)
    }()
    for {
        c.connected = true
        select {
        case message, ok := <-c.send:
            err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err != nil {
                loge.Error("Error calling conn.SetWriteDeadline error: %v", err)
            }

            if !ok {
                // The hub closed the channel.
                if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
                    return
                }
            }

            err = c.conn.WriteJSON(message)
            if err != nil {
                return
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

func (c *WSClient) handleMessage(message *ClientMessage) {

    if ! c.authenticated {

        if c.username == "" {

            if message.Payload != "" {
                c.authenticated = false
                c.username = message.Payload
                c.Send(SetAuthenticated(false))
                c.Send(SetEchoOn(false))
                c.Send(SetPrompt("Enter Password: "))
                return

            } else {
                c.authenticated = false
                c.Send(SetAuthenticated(false))
                c.Send(SetEchoOn(true))
                c.Send(SetPrompt("Enter Username: "))
                return
            }
        }

        if c.app.Config.UserAuthenticator(c, c.username, message.Payload) {
            c.app.Config.NotifyClientAuthenticated(c)
            c.authenticated = true

            if c.app.Config.DefaultPrompt != "" {
                c.Send(SetPrompt(c.app.Config.DefaultPrompt))
            }

            c.Send(SetEchoOn(true))
            c.Send(SetHistoryMode(false))
            c.Send(SetAuthenticated(true))
            return
        }

        c.authenticated = false
        c.loginAttempts++

        if c.loginAttempts < 3 {
            c.Send(SetAuthenticated(false))
            c.Send(SetEchoOn(false))
            c.Send(SetPrompt("Enter Password: "))
            c.Send(AppendText("Invalid Password", "red"))
            return
        }

        c.Send(SetAuthenticated(false))
        c.Send(SetPrompt(""))
        c.Send(AppendText("Invalid Password - disconnecting client", "red"))
        time.Sleep(500 * time.Millisecond)
        _ = c.conn.Close()

        return
    }

    if message.Payload == "exit" || message.Payload == "logoff" {
        c.Send(SetAuthenticated(false))
        c.Send(SetPrompt(""))
        c.Send(AppendText("logging off", "red"))

        time.Sleep(500 * time.Millisecond)
        _ = c.conn.Close()

    } else {
        c.history = append(c.history, message.Payload)
        c.app.Config.HandleCommand(c, message.Payload)
    }

}

// serveWs handles websocket requests from the peer.
func serveWs(app *MgmtApp, w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        loge.Error("serveWs error: %v\n", err)
        return
    }

    client := &WSClient{app: app, conn: conn, send: make(chan ServerMessage, 256), connected: true, httpReq: r}
    client.app.hub.register <- client
    client.misc = make(map[string]interface{})
    client.history = make([]string, 0)

    // Allow collection of memory referenced by the caller by doing all work in
    // new goroutines.
    go client.writePump()
    go client.readPump()

    app.Config.WelcomeUser(client)
}

func ConvertBytesToMessage(payload []byte) (*ClientMessage, error) {
    msg := &ClientMessage{}
    err := json.Unmarshal(payload, &msg)
    msg.Payload = strings.TrimSpace(msg.Payload)
    return msg, err
}

func (c *WSClient) IsAuthenticated() bool {
    return c.authenticated
}

func (c *WSClient) IsConnected() bool {
    return c.connected
}

func (c *WSClient) Username() string {
    return c.username
}

func (c *WSClient) History() []string {
    return c.history
}
func (c *WSClient) HttpReq() *http.Request {
    return c.httpReq
}
func (c *WSClient) Misc() map[string]interface{} {
    return c.misc
}

