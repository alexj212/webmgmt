package webmgmt

// BuildInfo is used to define the application build info, and inject values into via the build process.
type BuildInfo struct {
	// AppBuildInfo build information
	BuildDate string

	// LatestCommit build information
	LatestCommit string

	// BuildNumber build information
	BuildNumber string

	// BuiltOnIp build information
	BuiltOnIp string

	// BuiltOnOs build information
	BuiltOnOs string

	// RuntimeVer build information
	RuntimeVer string
}

var (
	// AppBuildInfo build information
	AppBuildInfo *BuildInfo
)
