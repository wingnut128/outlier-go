package version

// Version information, set at build time via ldflags
var (
	// Version is the semantic version (e.g., "1.0.0")
	Version = "dev"
	// GitCommit is the git commit hash
	GitCommit = "unknown"
	// BuildDate is the build timestamp
	BuildDate = "unknown"
)

// GetVersion returns the full version string
func GetVersion() string {
	return Version
}

// GetFullVersion returns version with git commit and build date
func GetFullVersion() string {
	if GitCommit != "unknown" {
		return Version + " (" + GitCommit[:7] + ", built " + BuildDate + ")"
	}
	return Version
}
