package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pkg/browser"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/potakhov/loge"
	"gitlab.paltalk.com/go/utils/netutils"

	"github.com/alexj212/webmgmt"
)

var (
	// BuildDate info from build
	BuildDate string

	// LatestCommit info from build
	LatestCommit string

	// BuildNumber info from build
	BuildNumber string

	// BuiltOnIp info from build
	BuiltOnIp string

	// BuiltOnOs info from build
	BuiltOnOs string

	// RuntimeVer info from build
	RuntimeVer string
)

var osSignal chan os.Signal
var onShutdownFunc func(os.Signal)
var logService *netutils.WsService
var httpListen = 1099
var listener net.Listener   // Listen socket for HTTP
var httpServer *http.Server // http server
var quit chan bool          // quit channel

func init() {
	osSignal = make(chan os.Signal, 1)
	onShutdownFunc = defaultShutdown
}

//nolint:gocyclo
func main() {
	var saveTemplateDir string
	flag.StringVar(&saveTemplateDir, "save", "", "save assets to directory")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] \n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	c := &customTransport{}

	logeShutdown := loge.Init(
		loge.Path("."),
		loge.EnableOutputConsole(true),
		loge.EnableOutputFile(false),
		loge.ConsoleOutput(os.Stdout),
		loge.EnableDebug(),
		loge.EnableError(),
		loge.EnableInfo(),
		loge.EnableWarning(),
		loge.Transports(func(list loge.TransactionList) []loge.Transport {
			transport := loge.WrapTransport(list, c)
			return []loge.Transport{transport}
		}),
	)

	defer logeShutdown()

	if saveTemplateDir != "" {
		err := webmgmt.SaveAssets(saveTemplateDir, webmgmt.DefaultWebEmbedFS, false)
		if err != nil {
			loge.Printf("Error writing assets: %v", err)
			os.Exit(-1)
		}
		os.Exit(1)
	}

	fmt.Printf("HttpListen           : %v\n", httpListen)

	loge.Info("NewMgmtApp: cfg.httpListen %v \n", httpListen)
	staticHtmlDir := "./web"

	root, err := webmgmt.SetupFS(webmgmt.DefaultWebEmbedFS, staticHtmlDir, false)
	if err != nil {
		loge.Error("error initializing fs: %v", err)
		os.Exit(-1)
	}

	listener, err = net.Listen("tcp", fmt.Sprintf(":%v", httpListen))
	if err != nil {
		loge.Error("error initializing http listener: %s", httpListen)

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
	app, err := setup(router, root)
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

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, app.WebPath(), 302)
	})

	router.NotFoundHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Printf("NotFoundHandler: %s\n", r.RequestURI)
		http.Redirect(rw, r, app.WebPath(), 302)
	})

	httpServer.Handler = router
	logService, err = loggerSetup(router)
	if err != nil {
		loge.Error("Error setting up loggerService error: ", err)
	}

	ticker := time.NewTicker(500 * time.Millisecond)
	quit := make(chan struct{})
	go func() {
		rand.Seed(time.Now().Unix())

		uids := []int{
			13564536,
			84,
			68450974,
			864567,
		}
		i := 0
		for {
			select {
			case <-ticker.C:
				uid := uids[rand.Intn(len(uids))]

				level := rand.Intn(6)
				switch level {
				case 0:
					loge.With("uid", uid).Debug("debug log message %d", i)
				case 1:
					loge.With("uid", uid).Info("info log message %d", i)
				case 2:
					loge.With("uid", uid).Trace("trace log message %d", i)
				case 3:
					loge.With("uid", uid).Warn("warn log message %d", i)
				case 4:
					loge.With("uid", uid).Error("error log message %d", i)
				default:
					loge.With("uid", uid).Printf("default log message %d", i)
				}

				i++
				// do stuff
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	LoopForever()

}

// LoopForever on signal processing
func LoopForever() {
	loge.Info("Entering infinite loop\n")

	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)
	sig := <-osSignal

	loge.Info("Exiting infinite loop received OsSignal\n")

	if onShutdownFunc != nil {
		onShutdownFunc(sig)
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

	url := fmt.Sprintf("http://localhost:%v", httpListen)
	browser.OpenURL(url)

	for {
		select {
		case <-quit:
			return nil

		case err := <-ch:
			return err
		}
	}
}

func shutdown() {
	loge.Info("MgmtApp Shutdown invoked")
	close(quit)
}

func (t *customTransport) WriteOutTransaction(tr *loge.Transaction) {
	fmt.Println("writeOutTransaction")
	//namespace string, room, event string, args ...interface{}
	for _, be := range tr.Items {

		payload, err := json.Marshal(be)
		if err != nil {
			loge.Error("Send json error: %v\n", err)
			continue
		}

		broadcaster.Send("logger", payload)
		//result := server.BroadcastToRoom("/", "logger", "message", be)
		//fmt.Printf("BroadcastToRoom: %v\n", result)
	}
}

func (t *customTransport) FlushTransactions() {
	fmt.Println("flushTransactions")
}

type customTransport struct {
}

var broadcaster netutils.Broadcast

// clientConn client connection struct
type clientConn struct {
	conn netutils.Connection
}

// createLoggerClient create connection for id. ID is assigned by service.
func createLoggerClient(conn netutils.Connection) (interface{}, error) {
	// var client ClientConn
	client := &clientConn{
		conn: conn,
	}
	return client, nil
}
func onClose(conn netutils.Connection) {
	broadcaster.Leave("logger", conn)

}
func onOpen(conn netutils.Connection) {
	broadcaster.Join("logger", conn)
}
func handle(conn netutils.Connection, payload []byte) {

}

func loggerSetup(router *mux.Router) (*netutils.WsService, error) {
	broadcaster = netutils.NewBroadcast()

	service, err := netutils.NewWsService(
		netutils.WsConnCreate(createLoggerClient),
		netutils.WsServiceOnClose(onClose),
		netutils.WsServiceOnOpen(onOpen),
		netutils.OnWsServiceRead(handle),
		netutils.WsServiceCheckOrigin(func(r *http.Request) bool { return true }),
	)

	if err != nil {
		loge.Printf("Error creating webTransport: %v\n", err)
		return nil, err
	}

	router.HandleFunc("/wslogger", func(w http.ResponseWriter, r *http.Request) {
		loge.Info("/wslogger invoked")
		service.ServeWs(w, r)
	})
	return service, nil
}
