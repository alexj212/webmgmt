package webmgmt

import (
    "net/http"
    "os"

    "github.com/gorilla/handlers"
    "github.com/gorilla/mux"
    "github.com/pkg/errors"
    "github.com/potakhov/loge"

    "github.com/gobuffalo/packr"
)

// http://localhost:7100/api/version
// http://localhost:7100/api/help
// http://localhost:7100/api/version/help
// http://localhost:7100/api/rooms/open/111
// http://localhost:7100/api/rooms/close/111
type BuildInfo struct {
    BuildDate    string
    LatestCommit string
    BuildNumber  string
    BuiltOnIp    string
    BuiltOnOs    string
    RuntimeVer   string
}

var (
    AppBuildInfo *BuildInfo
)

func (app *MgmtApp) initRouter(Name, InstanceId string) http.Handler {

    app.hub = newHub()
    go app.hub.run()

    upgrader.CheckOrigin = func(r *http.Request) bool { return true }

    loge.Info("app.Config.WebPath: %v\n", app.Config.WebPath)
    app.Config.Router.HandleFunc(app.Config.WebPath+"ws", func(w http.ResponseWriter, r *http.Request) {
        // loge.Info("/ws invoked")
        serveWs(app, w, r)
    })

    var fileHandler http.Handler

    fi, err := os.Stat("./web")
    if err == nil && fi.IsDir() {
        fileHandler = http.FileServer(http.Dir(app.Config.StaticHtmlDir))
    } else {
        box := packr.NewBox("./web")
        fileHandler = http.FileServer(box)
    }

    app.Config.Router.PathPrefix(app.Config.WebPath).Handler(http.StripPrefix(app.Config.WebPath, fileHandler))
    // app.router.PathPrefix("/").Handler(http.FileServer(http.Dir(app.Config.StaticHtmlDir)))

    reqHandlers := handlers.CORS(
        handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
        handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
        handlers.AllowedOrigins([]string{"*"}),
    )(app.Config.Router)

    // (Logger(os.Stderr, app.router))

    addHeaders := make(map[string]string)
    addHeaders["X-Served-By"] = "gooch1.0"
    addHeaders["X-Server-Id"] = InstanceId
    addHeaders["X-Server-Name"] = Name
    h := AddHeadersHandler(addHeaders, reqHandlers)

    return h
}

func (app *MgmtApp) handleServerVersion(w http.ResponseWriter, r *http.Request) {
    HttpNocacheJson(w)
    SendJson(w, r, AppBuildInfo)
}

type Results struct {
    Id        string      `json:"RoomId"`
    Action    string      `json:"Action"`
    Completed bool        `json:"Completed"`
    Data      interface{} `json:"data"`
}

func getRequestField(r *http.Request, fieldName string) (val string, err error) {
    valVar, ok := mux.Vars(r)[fieldName]
    if !ok {
        return "", errors.Errorf("Bad request: specify %v", fieldName)
    }

    return valVar, nil
}

type ErrorResult struct {
    Id    string `json:"id"`
    Error string `json:"error"`
}
