package main

import (
    "fmt"
    "io"
    "net"
    "net/http"
    "os"
    "os/signal"
    "runtime"
    "strings"
    "syscall"

    "github.com/alexj212/webmgmt"
    "github.com/gorilla/mux"
    "github.com/potakhov/loge"
)

const InstanceId = "InstanceId"
const Name = "Name"

var (
    BuildDate    string
    LatestCommit string
    BuildNumber  string
    BuiltOnIp    string
    BuiltOnOs    string
    RuntimeVer   string
)

var OsSignal chan os.Signal
var OnShutdownFunc func(os.Signal)

func init() {
    OsSignal = make(chan os.Signal, 1)
    OnShutdownFunc = defaultShutdown
}

var listener net.Listener   // Listen socket for HTTP
var httpServer *http.Server // http server
var quit chan bool          // quit channel

func main() {

    logeShutdown := loge.Init(
        loge.Path("."),
        loge.EnableOutputConsole(true),
        loge.EnableOutputFile(false),
        loge.ConsoleOutput(os.Stdout),
        loge.EnableDebug(),
        loge.EnableError(),
        loge.EnableInfo(),
        loge.EnableWarning(),
    )

    defer logeShutdown()

    HttpListen := 1099
    fmt.Printf("HttpListen           : %v\n", HttpListen)

    loge.Info("NewMgmtApp: cfg.httpListen %v \n", HttpListen)

    var err error
    listener, err = net.Listen("tcp", fmt.Sprintf(":%v", HttpListen))
    if err != nil {
        loge.Error("error initializing http listener: %s", HttpListen)

    }
    httpServer = &http.Server{}

    quit = make(chan bool)

    webmgmt.AppBuildInfo = &webmgmt.BuildInfo{}
    webmgmt.AppBuildInfo.BuildDate = BuildDate
    webmgmt.AppBuildInfo.LatestCommit = LatestCommit
    webmgmt.AppBuildInfo.BuildNumber = BuildNumber
    webmgmt.AppBuildInfo.BuiltOnIp = BuiltOnIp
    webmgmt.AppBuildInfo.BuiltOnOs = BuiltOnOs
    webmgmt.AppBuildInfo.RuntimeVer = runtime.Version()


    router := mux.NewRouter()
    _, err = Setup(router)
    if err != nil {
        loge.Error("Error starting server: %v", err)
        os.Exit(-1)
    }

    go func() {
        err := run()
        if err != nil {
            loge.Error("Event loop stopped with error: ", err)
        }
    }()



    router.HandleFunc("/hello/", func(w http.ResponseWriter, r *http.Request) {
        name := strings.Replace(r.URL.Path, "/hello/", "", 1)

        w.Header().Set("Content-Type", "text/plain")
        w.WriteHeader(http.StatusOK)

        io.WriteString(w, fmt.Sprintf("Hello %s\n", name))
    })

    router.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/plain")
        w.WriteHeader(http.StatusOK)

        io.WriteString(w, "Hello world\n")
    })

    router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/plain")
        w.WriteHeader(http.StatusOK)

        io.WriteString(w, "Hello world /\n")
    })

    httpServer.Handler = router
    LoopForever()

}

// Loop on signal processing
func LoopForever() {
    loge.Info("Entering infinite loop\n")

    signal.Notify(OsSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)
    sig := <-OsSignal

    loge.Info("Exiting infinite loop received OsSignal\n")

    if OnShutdownFunc != nil {
        OnShutdownFunc(sig)
    }
}

func defaultShutdown(sig os.Signal) {
    fmt.Printf("caught sig: %v\n\n", sig)
    os.Exit(0)
}

func run() error {
    loge.Info("EventLoop Run()")

    ch := make(chan error, 1)
    defer func() {
        err := listener.Close()
        if err != nil {
            loge.Error("Error closing listener error: %v", err)
        }

        err = httpServer.Close()
        if err != nil {
            loge.Error("Error closing http error: %v", err)
        }

    }()

    go func() {
        defer close(ch)
        loge.Info("Listening for HTTP on %v", listener.Addr())
        ch <- httpServer.Serve(listener)
    }()

    for {
        select {
        case <-quit:
            return nil

        case err := <-ch:
            return err
        }
    }
}

func Shutdown() {
    loge.Info("MgmtApp Shutdown invoked")
    close(quit)
}
