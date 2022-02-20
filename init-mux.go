package webmgmt

import (
    "github.com/gorilla/mux"
    "github.com/pkg/errors"
    "io/fs"
    "log"
    "net/http"
)

//InitMuxRouter will initialize the Router with the admin web app. It registers the webapp and assets file handler
// to be under the WebPath config field.
func InitMuxRouter(app *MgmtApp, router *mux.Router, webPath string, fileSystem fs.FS) error {

    upgrader.CheckOrigin = func(r *http.Request) bool { return true }

    // loge.Info("app.Config.WebPath: %v\n", app.webPath)
    router.HandleFunc(webPath+"ws", func(w http.ResponseWriter, r *http.Request) {
        // loge.Info("/ws invoked")
        serveWs(app, w, r)
    })

    router.HandleFunc(webPath+"version", app.handleServerVersion)

    if fileSystem != nil {

        var fileHandler http.Handler
        webDirHTTPFS := http.FS(fileSystem)
        fileHandler = http.FileServer(webDirHTTPFS)
        router.PathPrefix(webPath).Handler(http.StripPrefix(webPath, fileHandler))
    } else {
        log.Printf("unable to set file src or proxy\n")
        return errors.Errorf("unable to set file src or proxy")
    }

    return nil
}
