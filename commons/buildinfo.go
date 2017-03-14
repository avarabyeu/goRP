package commons

var (
	// Branch contains the current Git revision. Use make to build to make
	// sure this gets set.
	Branch string

	// BuildDate contains the date of the current build.
	BuildDate string

	// Version contains version
	Version string
)

//Build is global BuildInfo var
var Build *BuildInfo

// BuildInfo contains information about the current Hugo environment
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
