package webmgmt


// BuildInfo is used to define the application build info, and inject values into via the build process.
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
