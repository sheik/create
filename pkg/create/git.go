package create

func GitRepoClean() bool {
	return Command("git diff-index --quiet HEAD") == nil
}
