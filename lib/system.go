package lib

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/marguerite/util/command"
)

// FcCache run fc-cache command on the running system
func FcCache(verbosity int) {
	if cmd, err := command.Search("/usr/bin/fc-cache"); err == nil {
		Dbg(verbosity, Verbose, "Creating fontconfig cache files.\n")

		opts := ""

		if verbosity >= Verbose {
			opts = "--verbose"
		}

		_, status, _ := command.Run(cmd, opts)

		Dbg(verbosity, Debug, fmt.Sprintf("Exit status of fc-cache: %d\n", status))
	}
}

// FpRehash run xset fp rehash on the running system
func FpRehash(verbosity int) {
	if cmd, err := command.Search("/usr/bin/xset"); err == nil {
		re := regexp.MustCompile(`^:\d.*$`)
		if re.MatchString(GetEnv("DISPLAY")) {
			Dbg(verbosity, Verbose, "Rereading the font databases in the current font path ...\n")
			Dbg(verbosity, Debug, "Running xset fp rehash\n")

			out, _, _ := command.Run(cmd, "fp", "rehash")
			Dbg(verbosity, Debug, string(out)+"\n")
		} else {
			Dbg(verbosity, Verbose, "It is not a local display, do not reread X font databases for now.\n")
			Dbg(verbosity, Debug, "NOTE: do not run 'xset fp rehash', no local display detected.\n")
		}
	}
}

// ReloadXfsConfig reload Xorg Font Server on the running system
func ReloadXfsConfig(verbosity int) {
	if cmd, err := command.Search("/usr/bin/ps"); err == nil {
		pid, _, _ := command.Run(cmd, "-C", "xfs", "-o", "pid=")
		pid = strings.TrimSpace(pid)
		if len(pid) != 0 {
			Dbg(verbosity, Verbose, fmt.Sprintf("Reloading config file of X Font Server %s ...\n", pid))
			command.Run("/usr/bin/pkill", "-USR1", pid)
		} else {
			Dbg(verbosity, Debug, "X Font Server not used.\n")
		}
	} else {
		Dbg(verbosity, Verbose, "WARNING: ps command is missing, couldn't search for X Font Server pids.")
	}
}
