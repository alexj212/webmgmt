package main

import (
    "fmt"
    "os"
    "os/signal"
    "runtime"
    "syscall"
    "github.com/alexj212/webmgmt/webmgmt"

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

    config := &webmgmt.Config{HttpListen: 1099,StaticHtmlDir: "./web"}
    _, err := webmgmt.NewMgmtApp("testapp", "1", config)

    if err != nil {
        loge.Error("Error starting server: %v", err)
        os.Exit(-1)
    }
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
