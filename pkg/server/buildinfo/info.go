package buildinfo

var (
	// Version is the server version
	Version = "master"
	// CSSFiles is the css files
	CSSFiles = ""
	// JSFiles is the js files
	JSFiles = ""
	// RootURL is the root url
	RootURL = "/"
	// Standalone reprsents whether the build is for on-premises. It is a string
	// rather than a boolean, so that it can be overridden during compile time.
	Standalone = "false"
)
