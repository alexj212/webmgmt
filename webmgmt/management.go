package webmgmt

import (
    "fmt"
    "net"
    "net/http"
    "os"
    "sync"

    "github.com/gorilla/mux"
    "github.com/potakhov/loge"
)

type MgmtApp struct {
    Config     *Config
    quit       chan bool        // quit channel
    listener   net.Listener     // Listen socket for HTTP
    http       *http.Server     // http server
    router     *mux.Router      // http API router
    ch         chan interface{} // Main loop channel
    hub        *Hub
    opsCounter uint32
    ops        sync.Map
}

type Config struct {
    HttpListen                int
    StaticHtmlDir             string
    UserAuthenticator         func(client *Client, username string, password string) bool
    HandleCommand             func(c *Client, cmd string)
    NotifyClientAuthenticated func(client *Client)
    WelcomeUser               func(client *Client)
    UnregisterUser            func(client *Client)
    DefaultPrompt             string
}

func (cfg *Config) Display() {
    fmt.Println(os.Args)
    fmt.Println("-------------------------------------")
    fmt.Printf("HttpListen           : %v\n", cfg.HttpListen)
    fmt.Printf("StaticHtmlDir        : %s\n", cfg.StaticHtmlDir)
    fmt.Println("-------------------------------------")
}
func NewMgmtApp(name, instanceId string, config *Config) (*MgmtApp, error) {

    c := &MgmtApp{}
    c.Config = config

    if config.UserAuthenticator == nil {
        config.UserAuthenticator = c.userAuthenticator
    }
    if config.HandleCommand == nil {
        config.HandleCommand = c.handleCommand
    }
    if config.NotifyClientAuthenticated == nil {
        config.NotifyClientAuthenticated = c.notifyClientAuthenticated
    }
    if config.WelcomeUser == nil {
        config.WelcomeUser = c.welcomeUser
    }

    if config.UnregisterUser == nil {
        config.UnregisterUser = c.unregisterUser
    }

    if config.DefaultPrompt == "" {
        config.DefaultPrompt = "$"
    }

    c.quit = make(chan bool)
    c.ch = make(chan interface{}, 1)

    loge.Info("NewMgmtApp: cfg.httpListen %v \n", c.Config.HttpListen)

    l, err := net.Listen("tcp", fmt.Sprintf(":%v", c.Config.HttpListen))
    if err != nil {
        loge.Error("error initializing http listener: %s", c.Config.HttpListen)
        return nil, err
    }

    c.listener = l
    c.http = &http.Server{}
    c.initRouter(name, instanceId)

    go func() {
        err := c.run()
        if err != nil {
            loge.Error("Event loop stopped with error: ", err)
        }
    }()

    return c, nil
}

func (app *MgmtApp) Shutdown() {
    loge.Info("MgmtApp Shutdown invoked")
    close(app.quit)
}

func (app *MgmtApp) run() error {
    loge.Info("EventLoop Run()")

    ch := make(chan error, 1)
    defer func() {
        err := app.listener.Close()
        if err != nil {
            loge.Error("Error closing listener error: %v", err)
        }

        err = app.http.Close()
        if err != nil {
            loge.Error("Error closing http error: %v", err)
        }

    }()

    go func() {
        defer close(ch)
        loge.Info("Listening for HTTP on %v", app.listener.Addr())
        ch <- app.http.Serve(app.listener)
    }()

    for {
        select {
        case <-app.quit:
            return nil

        case err := <-ch:
            return err
        }
    }
}

func (app *MgmtApp) userAuthenticator(client *Client, username string, password string) bool {
    return true
}

func (app *MgmtApp) handleCommand(c *Client, cmd string) {
    c.Send(SetPrompt("$ "))
    c.Send(AppendText(fmt.Sprintf("echo: %v", cmd), "green"))
}

func (app *MgmtApp) notifyClientAuthenticated(client *Client) {
    loge.Info("New user on system: %v", client.Username())
}

func (app *MgmtApp) welcomeUser(client *Client) {
    client.Send(AppendText("Welcome to the machine", "red"))
    client.Send(SetEchoOn(true))
    client.Send(SetPrompt("Enter Username: "))
    client.Send(SetAuthenticated(false))
    client.Send(SetHistoryMode(false))
}

func (app *MgmtApp) unregisterUser(client *Client) {
    loge.Info("user logged off system: %v", client.Username())
}
