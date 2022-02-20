package webmgmt

import (
	"fmt"
	"github.com/alexj212/gox/ginx"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"io/fs"
	"log"
	"net/http"
	"strings"
)

//InitGin will initialize the Router with the admin web app. It registers the webapp and assets file handler
// to be under the WebPath config field.
func InitGin(app *MgmtApp, router *gin.Engine, webPath string, fileSystem fs.FS) error {

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// loge.Info("app.Config.WebPath: %v\n", app.webPath)
	router.GET(webPath+"ws", func(c *gin.Context) {
		// loge.Info("/ws invoked")
		serveWs(app, c.Writer, c.Request)
	})

	router.GET(webPath+"/version", func(c *gin.Context) {
		app.handleServerVersion(c.Writer, c.Request)
	})

	if fileSystem != nil {
		//newFS := &NewPathFS{
		//    base:   app.fileSystem,
		//    prefix: app.webPath,
		//}

		router.Use(ginx.Serve(webPath, ginx.StaticFS(fileSystem, webPath, true)))
	} else {
		log.Printf("unable to set file src or proxy\n")
		return errors.Errorf("unable to set file src or proxy")
	}
	return nil
}

func ServeProxy(app *MgmtApp, webPath string, filesProxy http.Handler) gin.HandlerFunc {
	handler := http.StripPrefix(webPath, filesProxy)

	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
		c.Abort()

	}
}

type NewPathFS struct {
	base   fs.FS
	prefix string
}

func (p *NewPathFS) Open(name string) (fs.File, error) {

	if !strings.HasPrefix(name, p.prefix) {
		fmt.Printf("ERROR Open: %v  prefix: %v\n", name, p.prefix)
		return nil, errors.Errorf("unable to open file: %v", name)
	}

	name = name[len(p.prefix):]
	fmt.Printf("Open: %v\n", name)
	return nil, nil
}
