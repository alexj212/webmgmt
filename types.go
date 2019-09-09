package webmgmt

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
