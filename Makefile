.PHONY: default Makefile

default: package

%:
	@go install github.com/sheik/create/cmd/create@latest
	@create $@

