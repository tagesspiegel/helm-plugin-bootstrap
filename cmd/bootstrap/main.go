package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/tagesspiegel/helm-plugin-bootstrap/internal/cli"
)

var (
	revision   string
	buildTime  time.Time
	dirtyBuild bool
	arch       string
	goos       string

	rootCmd = cli.NewBootstrapCmd(os.Stdout)
)

// we use the init function to setup the version information and flags for the root command
func init() {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	for _, set := range buildInfo.Settings {
		switch set.Key {
		case "vcs.revision":
			revision = set.Value
		case "vcs.time":
			buildTime, _ = time.Parse(time.RFC3339, set.Value)
		case "GOARCH":
			arch = set.Value
		case "GOOS":
			goos = set.Value
		case "vcs.modified":
			dirtyBuild = set.Value == "true"
		}
	}
	// add version info
	rootCmd.Version = fmt.Sprintf("%s %s - %v [ %s/%s ] [ %v ]", buildInfo.Main.Version, revision, buildTime, arch, goos, dirtyBuild)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
