package git

import (
	"errors"
	"github.com/sheik/create/pkg/shell"
)

func RepoClean(args ...interface{}) error {
	err := shell.Exec("git diff-index --quiet HEAD")
	if err != nil {
		return errors.New("git.RepoClean: git repository is dirty, commit and try again")
	}
	return nil
}

func CurrentTag() string {
	return shell.Output("git describe --tags | sed 's/-/_/g'")
}
