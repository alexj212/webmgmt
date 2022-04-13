package main

import (
    "embed"
    "fmt"
    "github.com/alexj212/gox"
    "github.com/alexj212/gox/httpx"
    "github.com/alexj212/webmgmt"
    "github.com/gin-gonic/gin"
    "github.com/gorilla/handlers"
    "github.com/gorilla/mux"
    "github.com/potakhov/loge"
    "io/fs"
    "net/http"
    "os"
)

var (
    //go:embed assets/*
    webEmbedFS embed.FS
    webFS      fs.FS
)

func launchMuxWeb(username, password, webPath string, httpPort int, webDir string) error {
    var err error
    webFS, err = fs.Sub(webEmbedFS, webDir)
    if err != nil {
        webFS = webmgmt.WebFS
    }

    root, err := gox.SetupFS(webFS, webDir, false)
    if err != nil {
        fmt.Printf("error initializing fs: %v\n", err)
        os.Exit(-1)
    }

    gox.WalkDir(root, "setupFS")

    app, err := initializeWebUi(username, password)
    if err != nil {
        fmt.Printf("Error starting server: %v\n", err)
        os.Exit(-1)
    }

    router := mux.NewRouter()

    //func InitMuxRouter(app *MgmtApp, router *mux.Router, webPath string, fileSystem fs.FS) error {
    err = webmgmt.InitMuxRouter(app, router, webPath, root)
    if err != nil {
        fmt.Printf("Error InitMuxRouter: %v\n", err)
        os.Exit(-1)
    }

    webmgmt.InitializeAuthMux("/auth", router)

    router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.Redirect(w, r, webPath, 302)
    })

    router.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf("hello \n")
        httpx.SendText(w, r, "hello world")
    })

    router.NotFoundHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
        fmt.Printf("NotFoundHandler: %s\n", r.RequestURI)
        http.Redirect(rw, r, webPath, 302)
    })

    reqHandlers := handlers.CORS(
        handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
        handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
        handlers.AllowedOrigins([]string{"*"}),
    )(router)

    addHeaders := make(map[string]string)
    addHeaders["X-Served-By"] = "gooch-1.0"
    addHeaders["X-Server-Id"] = "1"
    addHeaders["X-Server-Name"] = "example-mux"

    h := httpx.AddHeadersHandler(addHeaders, reqHandlers)

    go func() {
        fmt.Printf("HttpListen           : %v\n", httpPort)

        err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), h)
        if err != nil {
            loge.Error("Event loop stopped with error: %v", err)
            os.Exit(-1)
        }
    }()
    return nil
}

func launchGinWeb(username, password, webPath string, httpPort int, webDir string) (*gin.Engine, error) {
    var err error
    webFS, err = fs.Sub(webEmbedFS, webDir)
    if err != nil {
        webFS = webmgmt.WebFS
    }

    root, err := gox.SetupFS(webFS, webDir, false)
    if err != nil {
        fmt.Printf("error initializing fs: %v", err)
        os.Exit(-1)
    }

    gox.WalkDir(root, "setupFS")

    app, err := initializeWebUi(username, password)
    if err != nil {
        loge.Error("Error starting server: %v", err)
        os.Exit(-1)
    }

    router := gin.Default()

    webmgmt.InitializeAuthGin("/auth", router)

    router.Use(func() gin.HandlerFunc {
        return func(c *gin.Context) {
            c.Writer.Header().Set("X-Served-By", "gooch-1.0")
            c.Writer.Header().Set("X-Server-Id", "1")
            c.Writer.Header().Set("X-Server-Name", "example-mux")

            c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
            c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
            c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
            c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

            if c.Request.Method == "OPTIONS" {
                c.AbortWithStatus(204)
                return
            }

            c.Next()
        }
    }())

    err = webmgmt.InitGin(app, router, webPath, root)
    if err != nil {
        fmt.Printf("Error InitGin: %v\n", err)
        os.Exit(-1)
    }

    router.GET("/", func(c *gin.Context) {
        c.Redirect(http.StatusFound, webPath)
    })

    router.GET("/hello", func(c *gin.Context) {
        c.String(http.StatusOK, "hello world")
    })

    router.NoRoute(func(c *gin.Context) {
        c.Redirect(http.StatusFound, webPath)
    })

    go func() {
        fmt.Printf("HttpListen           : %v\n", httpPort)
        router.Run(fmt.Sprintf(":%d", httpPort))
    }()
    return router, nil

}
