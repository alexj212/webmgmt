package webmgmt

import (
    "fmt"
    "net/http"
    "os"

    "github.com/gorilla/mux"
    "github.com/pkg/errors"
    "github.com/potakhov/loge"
)

type MgmtApp struct {
    Config  *Config
    hub     *Hub
    Handler http.Handler
}

type Config struct {
    StaticHtmlDir             string
    UserAuthenticator         func(client Client, username string, password string) bool
    HandleCommand             func(c Client, cmd string)
    NotifyClientAuthenticated func(client Client)
    WelcomeUser               func(client Client)
    UnregisterUser            func(client Client)
    DefaultPrompt             string
    Router                    *mux.Router // http API router
    WebPath                   string
}

func (cfg *Config) Display() {
    fmt.Println(os.Args)
    fmt.Println("-------------------------------------")

    fmt.Printf("StaticHtmlDir        : %s\n", cfg.StaticHtmlDir)
    fmt.Println("-------------------------------------")
}
func NewMgmtApp(name, instanceId string, config *Config) (*MgmtApp, error) {

    if config == nil {
        return nil, errors.Errorf("cannot create MgmtApp with nil config")
    }

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

    c.Handler = c.initRouter(name, instanceId)
    return c, nil
}

func (app *MgmtApp) userAuthenticator(client Client, username string, password string) bool {
    return true
}

func (app *MgmtApp) handleCommand(c Client, cmd string) {
    c.Send(SetPrompt("$ "))
    c.Send(AppendText(fmt.Sprintf("echo: %v", cmd), "green"))
}

func (app *MgmtApp) notifyClientAuthenticated(client Client) {
    loge.Info("New user on system: %v", client.Username())
}

func (app *MgmtApp) welcomeUser(client Client) {
    client.Send(AppendText("Welcome to the machine", "red"))
    client.Send(SetEchoOn(true))
    client.Send(SetPrompt("Enter Username: "))
    client.Send(SetAuthenticated(false))
    client.Send(SetHistoryMode(false))
}

func (app *MgmtApp) unregisterUser(client Client) {
    loge.Info("user logged off system: %v", client.Username())
}
