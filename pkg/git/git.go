package git

import (
	"errors"
	"fmt"
	"github.com/sheik/create/pkg/shell"
	"strconv"
	"strings"
)

func RepoClean(args ...interface{}) error {
	err := shell.Exec("git diff-index --quiet HEAD")
	if err != nil {
		return errors.New("git.RepoClean: git repository is dirty, commit and try again")
	}
	return nil
}
func IncrementMinorVersion(version string) string {
	parts := strings.Split(strings.Split(version, "_")[0], ".")
	if len(parts) != 3 {
		return version
	}
	if minorVersion, err := strconv.Atoi(parts[2]); err == nil {
		minorVersion += 1
		return fmt.Sprintf("%s.%s.%d", parts[0], parts[1], minorVersion)
	} else {
		return version
	}
}

func GetVersion() string {
	return shell.Output("git describe --tags | sed 's/-/_/g'")
}
