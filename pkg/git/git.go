package git

import "github.com/sheik/create/pkg/create"

func RepoClean() bool {
	return create.Command("git diff-index --quiet HEAD") == nil
}

func CurrentTag() string {
	return create.Output("git describe --tags | sed 's/-/_/g'")
}
