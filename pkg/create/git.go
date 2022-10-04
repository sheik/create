package create

func GitRepoDirty() bool {
	return Command("git diff-index --quiet HEAD") == nil
}
