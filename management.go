package webmgmt

import (
	"fmt"
	"os"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/potakhov/loge"
)

// MgmtApp struct is the web admin app.
type MgmtApp struct {
	staticHtmlDir                   string
	userAuthenticator               func(client Client, username string, password string) bool
	handleCommand                   func(c Client, cmd string)
	notifyClientAuthenticated       func(client Client)
	notifyClientAuthenticatedFailed func(client Client)
	welcomeUser                     func(client Client)
	unregisterUser                  func(client Client)
	defaultPrompt                   string
	webPath                         string
	clientInitializer               func(client Client)
	hub                             *Hub
}

// Config struct  is used to configure a WebMgmt admin handler.
type Config struct {
	StaticHtmlDir                   string
	DefaultPrompt                   string
	WebPath                         string
	UserAuthenticator               func(client Client, username string, password string) bool
	HandleCommand                   func(c Client, cmd string)
	NotifyClientAuthenticated       func(client Client)
	notifyClientAuthenticatedFailed func(client Client)
	WelcomeUser                     func(client Client)
	UnregisterUser                  func(client Client)
	ClientInitializer               func(client Client)
}

// Display is used to display the config
func (cfg *Config) Display() {
	fmt.Println(os.Args)
	fmt.Println("-------------------------------------")
	fmt.Printf("StaticHtmlDir        : %s\n", cfg.StaticHtmlDir)
	fmt.Println("-------------------------------------")
}

// NewMgmtApp will create a new web mgmt web handler with the Config passed in, Various funcs can be overwritten for authentication, welcome etc. If an error is encountered,
// it will be returned
func NewMgmtApp(name, instanceId string, config *Config, router *mux.Router) (*MgmtApp, error) {

	if config == nil {
		return nil, errors.Errorf("cannot create MgmtApp with nil config")
	}

	c := &MgmtApp{}
	c.staticHtmlDir = config.StaticHtmlDir
	c.userAuthenticator = config.UserAuthenticator
	c.handleCommand = config.HandleCommand
	c.notifyClientAuthenticated = config.NotifyClientAuthenticated
	c.notifyClientAuthenticatedFailed = config.notifyClientAuthenticatedFailed
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
	if c.notifyClientAuthenticatedFailed == nil {
		c.notifyClientAuthenticatedFailed = c.defaultNotifyClientAuthenticatedFailed
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

// defaultUserAuthenticator the default function that is invoked when a client is attempting to login. The username and password is passed to the func to be validated.
// It can be overwritten via a func in the initialization Config
func (app *MgmtApp) defaultUserAuthenticator(client Client, username string, password string) bool {
	return true
}

// defaultHandleCommand the default function that is invoked when a client sends text to the server.  It can be overwritten via a func in the initialization Config
func (app *MgmtApp) defaultHandleCommand(c Client, cmd string) {
	c.Send(SetPrompt("$ "))
	c.Send(AppendText(fmt.Sprintf("echo: %v", cmd), "green"))
}

// defaultNotifyClientAuthenticated the default function that is invoked when a client is authenticated on the system.  It can be overwritten via a func in the initialization Config
func (app *MgmtApp) defaultNotifyClientAuthenticated(client Client) {
	loge.Info("New user on system: %v", client.Username())
}

// defaultNotifyClientAuthenticatedFailed the default function that is invoked when a client fails authenticated on the system.  It can be overwritten via a func in the initialization Config
func (app *MgmtApp) defaultNotifyClientAuthenticatedFailed(client Client) {
	loge.Info("User failed login: %v", client.Username())
}

// defaultWelcomeUser the default function which is invoked to welcome a user. It can be overwritten via a func in the initialization Config
func (app *MgmtApp) defaultWelcomeUser(client Client) {
	client.Send(AppendText("Welcome to the machine", "red"))
	client.Send(SetEchoOn(true))
	client.Send(SetPrompt("Enter Username: "))
	client.Send(SetAuthenticated(false))
	client.Send(SetHistoryMode(false))
}

// defaultUnregisterUser is the default client unregister function.. It can be overwritten by setting a new func in the Config. It is used
// when a client is disconnected from the server, the function will be invoked with the Client structure.
func (app *MgmtApp) defaultUnregisterUser(client Client) {
	loge.Info("user logged off system: %v", client.Username())
}

// defaultClientInitializer is the default client initializer. It can be overwritten by setting a new func in the Config. If it used
// when a new client is connected, the function will be invoked with the Client structure.
func (app *MgmtApp) defaultClientInitializer(client Client) {
	loge.Info("Default defaultClientInitializer: %v", client.Username())
}

// Broadcast will take a ServerMessage and send to all clients that are currently connected.
func (app *MgmtApp) Broadcast(msg ServerMessage) {
	app.hub.Broadcast(msg)
}

// WebPath will return the web path that has been defined for the MgmtAoo web app.
func (app *MgmtApp) WebPath() string {
	return app.webPath
}
