package main

import (
    "fmt"
    "github.com/alexj212/gox"
    "github.com/alexj212/gox/utilx"
    "github.com/alexj212/webmgmt"
    "os"
    "runtime"

    "github.com/droundy/goopt"
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

func init() {
    webmgmt.AppBuildInfo = &webmgmt.BuildInfo{}
    webmgmt.AppBuildInfo.BuildDate = BuildDate
    webmgmt.AppBuildInfo.LatestCommit = LatestCommit
    webmgmt.AppBuildInfo.BuildNumber = BuildNumber
    webmgmt.AppBuildInfo.BuiltOnIp = BuiltOnIp
    webmgmt.AppBuildInfo.BuiltOnOs = BuiltOnOs
    webmgmt.AppBuildInfo.RuntimeVer = runtime.Version()

    goopt.Description = func() string {
        return "web mgt test app"
    }
    goopt.Author = "Alex Jeannopoulos"
    goopt.ExtraUsage = ``
    goopt.Summary = `
        Simple usage test application.
`

    goopt.Version = fmt.Sprintf(
        `build information

  LatestCommit  : %s
  BuildNumber   : %s
  BuiltOnIp     : %s
  BuiltOnOs     : %s
  RuntimeVer    : %s
  BuildDate     : %s
`, LatestCommit, BuildNumber, BuiltOnIp, BuiltOnOs, RuntimeVer, BuildDate)

    //Parse options
    goopt.Parse(nil)

}

var (
    exportTemplates = goopt.Flag([]string{"--export"}, nil, "export templates to --webDir value.", "")
    webDir          = goopt.String([]string{"--webDir"}, "./assets", "web assets directory")
    adminWebPath    = goopt.String([]string{"--path"}, "/admin/", "admin web path")
    httpPort        = goopt.Int([]string{"--port"}, 1099, "port for server")
    adminUserName   = goopt.String([]string{"--username"}, "admin", "admin username")
    adminPassword   = goopt.String([]string{"--password"}, "bambam", "admin password")

    useGin = goopt.Flag([]string{"--gin"}, nil, "use gin router.", "")
    useMux = goopt.Flag([]string{"--mux"}, nil, "use mux router.", "")
)

func main() {

    if *exportTemplates {
        err := gox.SaveAssets(*webDir, webmgmt.DefaultWebEmbedFS, false)
        if err != nil {
            fmt.Printf("Error writing assets: %v", err)
            os.Exit(-1)
        }
        os.Exit(1)
    }

    if *useGin && *useMux {
        fmt.Printf("You cannot specify, --gin and --mux\n")
        os.Exit(-1)
    }

    var err error

    _, err = launchGinWeb(*adminUserName, *adminPassword, *adminWebPath, *httpPort, *webDir)
    if err != nil {
        fmt.Printf("Error starting admin server: %v", err)
        os.Exit(-1)
    }

    err = launchMuxWeb(*adminUserName, *adminPassword, *adminWebPath, *httpPort+1, *webDir)
    if err != nil {
        fmt.Printf("Error starting admin server: %v", err)
        os.Exit(-1)
    }

    utilx.LoopForever(nil)
}
