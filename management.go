package webmgmt

import (
    "fmt"
    "os"

    "github.com/gorilla/mux"
    "github.com/pkg/errors"
    "github.com/potakhov/loge"
)

type MgmtApp struct {
    staticHtmlDir             string
    userAuthenticator         func(client Client, username string, password string) bool
    handleCommand             func(c Client, cmd string)
    notifyClientAuthenticated func(client Client)
    welcomeUser               func(client Client)
    unregisterUser            func(client Client)
    defaultPrompt             string
    webPath                   string
    clientInitializer           func(client Client)
    hub                       *Hub
}

type Config struct {
    StaticHtmlDir             string
    DefaultPrompt             string
    WebPath                   string
    UserAuthenticator         func(client Client, username string, password string) bool
    HandleCommand             func(c Client, cmd string)
    NotifyClientAuthenticated func(client Client)
    WelcomeUser               func(client Client)
    UnregisterUser            func(client Client)
    ClientInitializer           func(client Client)
}

func (cfg *Config) Display() {
    fmt.Println(os.Args)
    fmt.Println("-------------------------------------")
    fmt.Printf("StaticHtmlDir        : %s\n", cfg.StaticHtmlDir)
    fmt.Println("-------------------------------------")
}
func NewMgmtApp(name, instanceId string, config *Config, router *mux.Router) (*MgmtApp, error) {

    if config == nil {
        return nil, errors.Errorf("cannot create MgmtApp with nil config")
    }

    c := &MgmtApp{}
    c.staticHtmlDir = config.StaticHtmlDir
    c.userAuthenticator = config.UserAuthenticator
    c.handleCommand = config.HandleCommand
    c.notifyClientAuthenticated = config.NotifyClientAuthenticated
    c.welcomeUser = config.WelcomeUser
    c.unregisterUser = config.UnregisterUser
    c.defaultPrompt = config.DefaultPrompt
    c.webPath = config.WebPath
    c.clientInitializer = config.ClientInitializer

    if c.userAuthenticator == nil {
        c.userAuthenticator = c.defaultUserAuthenticator
    }
    if c.handleCommand == nil {
        c.handleCommand = c.defaultHandleCommand
    }
    if c.notifyClientAuthenticated == nil {
        c.notifyClientAuthenticated = c.defaultNotifyClientAuthenticated
    }
    if c.welcomeUser == nil {
        c.welcomeUser = c.defaultWelcomeUser
    }

    if c.unregisterUser == nil {
        c.unregisterUser = c.defaultUnregisterUser
    }

    if c.clientInitializer == nil {
        c.clientInitializer = c.defaultClientInitializer
    }

    if c.defaultPrompt == "" {
        c.defaultPrompt = "$"
    }

    c.initRouter(name, instanceId, router)
    return c, nil
}

func (app *MgmtApp) defaultUserAuthenticator(client Client, username string, password string) bool {
    return true
}

func (app *MgmtApp) defaultHandleCommand(c Client, cmd string) {
    c.Send(SetPrompt("$ "))
    c.Send(AppendText(fmt.Sprintf("echo: %v", cmd), "green"))
}

func (app *MgmtApp) defaultNotifyClientAuthenticated(client Client) {
    loge.Info("New user on system: %v", client.Username())
}

func (app *MgmtApp) defaultWelcomeUser(client Client) {
    client.Send(AppendText("Welcome to the machine", "red"))
    client.Send(SetEchoOn(true))
    client.Send(SetPrompt("Enter Username: "))
    client.Send(SetAuthenticated(false))
    client.Send(SetHistoryMode(false))
}

func (app *MgmtApp) defaultUnregisterUser(client Client) {
    loge.Info("user logged off system: %v", client.Username())
}

func (app *MgmtApp) defaultClientInitializer(client Client) {
    loge.Info("Default defaultClientInitializer: %v", client.Username())
}

func (app *MgmtApp) Broadcast(msg ServerMessage) {
    app.hub.Broadcast(msg)
}
