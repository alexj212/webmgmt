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
    Config *Config

    quit                   chan bool        // quit channel
    listener               net.Listener     // Listen socket for HTTP
    http                   *http.Server     // http server
    router                 *mux.Router      // http API router
    ch                     chan interface{} // Main loop channel
    hub                    *Hub
    opsCounter uint32
    ops        sync.Map
}


type Config struct {
    HttpListen      int
    StaticHtmlDir   string
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

    loge.Info("NewMgmtApp: cfg.httpListen %v \n", c.Config.HttpListen)



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
