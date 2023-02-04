package build

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"

	log "github.com/sirupsen/logrus"
)

var GOOS string = runtime.GOOS
var GOARCH string = runtime.GOARCH
var RuntimeVersion string = strings.TrimLeft(runtime.Version(), "go")
var VcsRevision string
var VcsTag string
var VcsTime string
var Version string

func init() {
	info, ok := debug.ReadBuildInfo()

	if ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" && VcsRevision == "" {
				VcsRevision = setting.Value
			} else if setting.Key == "vcs.time" && VcsTime == "" {
				VcsTime = setting.Value
			}
		}
	}
}

func GetVersion() string {
	var version string = fmt.Sprintf("%s/%s %s/%s", Version, VcsTag, GOOS, GOARCH)

	if log.GetLevel() == log.DebugLevel {
		fields := log.Fields{
			"version":  fmt.Sprintf("%s/%s", Version, VcsTag),
			"platform": fmt.Sprintf("%s/%s", GOOS, GOARCH),
			"date":     VcsTime,
			"go":       RuntimeVersion,
		}

		log.WithFields(fields).Debugf("%s %s go%s", version, VcsTime, RuntimeVersion)
	}

	return version
}
