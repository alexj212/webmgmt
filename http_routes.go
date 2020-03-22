package webmgmt

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gobuffalo/packd"

	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

//initRouter will initialize the Router with the admin web app. It registers the webapp and assets file handler
// to be under the WebPath config field.
func (app *MgmtApp) initRouter(Name, InstanceId string, router *mux.Router) http.Handler {

	app.hub = newHub()
	go app.hub.run()

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// loge.Info("app.Config.WebPath: %v\n", app.webPath)
	router.HandleFunc(app.webPath+"ws", func(w http.ResponseWriter, r *http.Request) {
		// loge.Info("/ws invoked")
		serveWs(app, w, r)
	})

	var fileHandler http.Handler

	fileHandler = http.FileServer(app.fileSystem)
	router.HandleFunc(app.webPath+"/version", app.handleServerVersion)
	router.PathPrefix(app.webPath).Handler(http.StripPrefix(app.webPath, fileHandler))

	reqHandlers := handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}),
	)(router)

	// (Logger(os.Stderr, app.router))

	addHeaders := make(map[string]string)
	addHeaders["X-Served-By"] = "gooch-1.0"
	addHeaders["X-Server-Id"] = InstanceId
	addHeaders["X-Server-Name"] = Name
	h := AddHeadersHandler(addHeaders, reqHandlers)

	return h
}

// handleServerVersion handles the server version rest handler
func (app *MgmtApp) handleServerVersion(w http.ResponseWriter, r *http.Request) {
	HttpNocacheJson(w)
	SendJson(w, r, AppBuildInfo)
}

// SaveTemplates will save the prepacked templates for local editing. File structure will be recreated under the output dir.
func SaveAssets(outputDir string) error {
	fmt.Printf("SaveAssets: %v\n", outputDir)
	if outputDir == "" {
		outputDir = "."
	}

	if strings.HasSuffix(outputDir, "/") {
		outputDir = outputDir[:len(outputDir)-1]
	}

	if outputDir == "" {
		outputDir = "."
	}

	box := packr.New("webmgmt", "./web")

	box.Walk(func(s string, file packd.File) error {
		fileName := fmt.Sprintf("%s/%s", outputDir, s)

		fi, err := file.FileInfo()
		if err == nil {
			if !fi.IsDir() {

				err := WriteNewFile(fileName, file)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	return nil
}
