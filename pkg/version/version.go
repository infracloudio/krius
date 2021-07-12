package version

import "fmt"

// The below variables are overrriden using the build process
var Version = "dev"
var GitCommitID = "none"
var BuildDate = "unknown"

const versionLongFmt = `{"Version": "%s", "GitCommit": "%s", "BuildDate": "%s"}`

func Long() string {
	return fmt.Sprintf(versionLongFmt, Version, GitCommitID, BuildDate)
}

func Short() string {
	return Version
}
