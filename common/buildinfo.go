package common

var (
	// Branch contains the current Git revision. Use make to build to make
	// sure this gets set.
	Branch string

	// BuildDate contains the date of the current build.
	BuildDate string

	// Version contains version
	Version string
)

var Build *BuildInfo

// HugoInfo contains information about the current Hugo environment
type BuildInfo struct {
	Version   string `json:"version,omitempty"`
	Branch    string `json:"branch,omitempty"`
	BuildDate string `json:"build_date,omitempty"`
}

func init() {
	Build = &BuildInfo{
		Version:   Version,
		Branch:    Branch,
		BuildDate: BuildDate,
	}
}
