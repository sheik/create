package git

import (
	"github.com/sheik/create/pkg/shell"
)

func RepoClean() bool {
	return shell.Exec("git diff-index --quiet HEAD") == nil
}

func CurrentTag() string {
	return shell.Output("git describe --tags | sed 's/-/_/g'")
}
