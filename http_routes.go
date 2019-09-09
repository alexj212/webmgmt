package webmgmt

import (
    "net/http"
    "os"

    "github.com/gorilla/handlers"
    "github.com/gorilla/mux"
    "github.com/potakhov/loge"

    "github.com/gobuffalo/packr"
)



func (app *MgmtApp) initRouter(Name, InstanceId string, router *mux.Router) http.Handler {

    app.hub = newHub()
    go app.hub.run()

    upgrader.CheckOrigin = func(r *http.Request) bool { return true }

    loge.Info("app.Config.WebPath: %v\n", app.webPath)
    router.HandleFunc(app.webPath+"ws", func(w http.ResponseWriter, r *http.Request) {
        // loge.Info("/ws invoked")
        serveWs(app, w, r)
    })

    var fileHandler http.Handler

    fi, err := os.Stat("./web")
    if err == nil && fi.IsDir() {
        fileHandler = http.FileServer(http.Dir(app.staticHtmlDir))
    } else {
        box := packr.NewBox("./web")
        fileHandler = http.FileServer(box)
    }

    router.PathPrefix(app.webPath).Handler(http.StripPrefix(app.webPath, fileHandler))
    // app.router.PathPrefix("/").Handler(http.FileServer(http.Dir(app.Config.StaticHtmlDir)))

    reqHandlers := handlers.CORS(
        handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
        handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
        handlers.AllowedOrigins([]string{"*"}),
    )(router)

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

type ErrorResult struct {
    Id    string `json:"id"`
    Error string `json:"error"`
}
