package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
	"github.com/potakhov/loge"

	"github.com/alexj212/webmgmt"
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
	var saveTemplateDir string
	flag.StringVar(&saveTemplateDir, "save", "", "save assets to directory")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] \n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

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

	var root http.FileSystem
	staticHtmlDir := "./web"

	if staticHtmlDir != "" {
		fi, err := os.Stat(staticHtmlDir)
		if err == nil && fi.IsDir() {
			loge.Info("using file serving from local disk: %v\n", fi.Name())
			root = http.Dir(staticHtmlDir)
		}
	}

	if root == nil {
		loge.Info("using file serving from packed resources \n")
		box := packr.New("webmgmt", staticHtmlDir)
		loge.Info("Box Details: Name: %v Path: %v ResolutionDir: %v", box.Name, box.Path, box.ResolutionDir)

		for i, file := range box.List() {
			loge.Info("   [%d] [%s]", i, file)
		}
		root = box
	}

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
	app, err := Setup(router, root)
	if err != nil {
		loge.Error("Error starting server: %v", err)
		os.Exit(-1)
	}

	if saveTemplateDir != "" {
		err = webmgmt.SaveAssets(saveTemplateDir)
		if err != nil {
			loge.Printf("Error writing assets: %v", err)
			os.Exit(-1)
		}
	}

	go func() {
		err := run()
		if err != nil {
			loge.Error("Event loop stopped with error: ", err)
		}
	}()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, app.WebPath(), 302)
	})

	router.NotFoundHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		http.Redirect(rw, r, app.WebPath(), 302)
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
