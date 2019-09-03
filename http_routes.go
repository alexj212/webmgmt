package webmgmt

import (
    "net/http"

    "github.com/gorilla/handlers"
    "github.com/gorilla/mux"
    "github.com/pkg/errors"
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

func (app *MgmtApp) initRouter(Name, InstanceId string) {
    app.router = mux.NewRouter()

    app.hub = newHub()
    go app.hub.run()

    upgrader.CheckOrigin = func(r *http.Request) bool { return true }

    app.router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        // loge.Info("/ws invoked")
        serveWs(app, w, r)
    })

    app.router.PathPrefix("/").Handler(http.FileServer(http.Dir(app.Config.StaticHtmlDir)))

    reqHandlers := handlers.CORS(
        handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
        handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
        handlers.AllowedOrigins([]string{"*"}),
    )(app.router)

    // (Logger(os.Stderr, app.router))

    addHeaders := make(map[string]string)
    addHeaders["X-Served-By"] = "gooch1.0"
    addHeaders["X-Server-Id"] = InstanceId
    addHeaders["X-Server-Name"] = Name
    h := AddHeadersHandler(addHeaders, reqHandlers)

    app.http.Handler = h
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
