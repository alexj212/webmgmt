package main

import (
    "fmt"
    "os"
    "os/signal"
    "runtime"
    "syscall"
    "time"

    "github.com/alexj212/webmgmt"
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

    webmgmt.AppBuildInfo = &webmgmt.BuildInfo{}
    webmgmt.AppBuildInfo.BuildDate = BuildDate
    webmgmt.AppBuildInfo.LatestCommit = LatestCommit
    webmgmt.AppBuildInfo.BuildNumber = BuildNumber
    webmgmt.AppBuildInfo.BuiltOnIp = BuiltOnIp
    webmgmt.AppBuildInfo.BuiltOnOs = BuiltOnOs
    webmgmt.AppBuildInfo.RuntimeVer = runtime.Version()

    config := &webmgmt.Config{HttpListen: 1099, StaticHtmlDir: "./web"}
    config.DefaultPrompt = "$"

    config.UserAuthenticator = func(client *webmgmt.Client, s string, s2 string) bool {
        return s == "alex" && s2 == "bambam"
    }

    commands := []string{"http", "user", "prompt", "link", "ticker", "image", "raw", "commands", "history"}

    config.HandleCommand = func(client *webmgmt.Client, cmd string) {
        // loge.Info("handleMessage  - authenticated user message.Payload: [" + cmd+"]")

        switch cmd {

        case "ticker":
            toggleTicker(client)

        case "image":
            client.Send(webmgmt.AppendRawText("<img width=\"200\" height=\"200\" src=\"https://avatars1.githubusercontent.com/u/174203?s=200&v=4\" alt=\"me\"/>", nil))

        case "raw":
            client.Send(webmgmt.AppendRawText("<img width=\"100\" height=\"100\" src=\"https://avatars1.githubusercontent.com/u/174203?s=100&v=4\" alt=\"me\"/>", commands))

        case "commands":
            client.Send(webmgmt.AppendRawText("Set Commands", commands))

        case "link":
            client.Send(webmgmt.AppendRawText(webmgmt.Link("http://www.slashdot.org", webmgmt.Color("orange", "slashdot")), nil))

        case "prompt":
            client.Send(webmgmt.SetPrompt(webmgmt.Color("red", client.Username()) + "@" + webmgmt.Color("green", "myserver") + ":&nbsp;"))

        case "http":
            displayHttpInfo(client)

        case "history":
            displayHistory(client)

        case "user":
            displayUserInfo(client)

        default:
            client.Send(webmgmt.AppendText(fmt.Sprintf("echo: %v", cmd), "green"))
        }
    }

    config.NotifyClientAuthenticated = func(client *webmgmt.Client) {
        client.Send(webmgmt.SetPrompt("$ "))
        loge.Info("New user authenticated on system: %v", client.Username())
    }

    config.UnregisterUser = func(client *webmgmt.Client) {
        loge.Info("user logged off system: %v", client.Username())
    }

    config.WelcomeUser = func(client *webmgmt.Client) {
        client.Send(webmgmt.AppendText("Welcome to the machine", "red"))
        client.Send(webmgmt.SetPrompt("Enter Username: "))
        client.Send(webmgmt.SetAuthenticated(false))
        client.Send(webmgmt.SetHistoryMode(false))
        client.Send(webmgmt.SetEchoOn(true))

    }

    _, err := webmgmt.NewMgmtApp("testapp", "1", config)

    if err != nil {
        loge.Error("Error starting server: %v", err)
        os.Exit(-1)
    }
    LoopForever()

}

func displayHistory(client *webmgmt.Client) {
    if len(client.History) > 0 {
        for i, cmd := range client.History {
            client.Send(webmgmt.AppendText(fmt.Sprintf("History[%d]: %v", i, cmd), "green"))
        }
    } else {
        client.Send(webmgmt.AppendText("History is empty", "green"))
    }
}

func displayUserInfo(client *webmgmt.Client) {
    client.Send(webmgmt.AppendText(fmt.Sprintf("Username       : %v", client.Username()), "green"))
    client.Send(webmgmt.AppendText(fmt.Sprintf("IsAuthenticated: %v", client.IsAuthenticated()), "green"))
}

func displayHttpInfo(client *webmgmt.Client) {
    client.Send(webmgmt.AppendText(fmt.Sprintf("GetIPAdress              : %v", webmgmt.GetIPAddress(client.HttpReq)), "green"))
    client.Send(webmgmt.AppendText(fmt.Sprintf("client.HttpReq.Host      : %v", client.HttpReq.Host), "green"))
    client.Send(webmgmt.AppendText(fmt.Sprintf("client.HttpReq.Method    : %v", client.HttpReq.Method), "green"))
    client.Send(webmgmt.AppendText(fmt.Sprintf("client.HttpReq.RemoteAddr: %v", client.HttpReq.RemoteAddr), "green"))
    client.Send(webmgmt.AppendText(fmt.Sprintf("client.HttpReq.RequestURI: %v", client.HttpReq.RequestURI), "green"))
    client.Send(webmgmt.AppendText(fmt.Sprintf("client.HttpReq.Referer() : %v", client.HttpReq.Referer()), "green"))
    for i, cookie := range client.HttpReq.Cookies() {
        client.Send(webmgmt.AppendText(fmt.Sprintf("client.HttpReq.Cookies[%-2d]: %-25v / %v", i, cookie.Name, cookie.Value), "green"))
    }

    for name, values := range client.HttpReq.Header {
        // Loop over all values for the name.
        for _, value := range values {
            client.Send(webmgmt.AppendText(fmt.Sprintf("client.HttpReq.Header[%-25v]:  %v", name, value), "green"))
        }
    }

}

func toggleTicker(client *webmgmt.Client) {

    tickerDone, ok := client.Misc["ticker_done"]
    if ok {
        delete(client.Misc, "ticker_done")

        var done chan bool
        done = tickerDone.(chan bool)
        done <- true
    } else {
        done := make(chan bool)
        client.Misc["ticker_done"] = done

        go func() {
            ticker := time.NewTicker(5 * time.Second)

            for {
                if !client.IsConnected() {
                    break
                }

                select {

                case <-done:
                    loge.Info("Ticker Done")
                    return

                case t := <-ticker.C:
                    loge.Info("Ticker ticked")
                    if client.IsAuthenticated() {
                        msg := fmt.Sprintf("Tick at: %v", t)
                        client.Send(webmgmt.AppendText(msg, "blue"))
                    }
                }
            }
            loge.Info("Ticker Stopped")
            ticker.Stop()
        }()

    }

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
