package git

import (
	"github.com/sheik/create/pkg/create"
	"github.com/sheik/create/pkg/shell"
)

func RepoClean() bool {
	return shell.Exec("git diff-index --quiet HEAD") == nil
}

func CurrentTag() string {
	return create.Output("git describe --tags | sed 's/-/_/g'")
}
