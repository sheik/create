.PHONY: default install Makefile

default: package

install:
	@go install github.com/sheik/create/cmd/create@latest

%: install
	@create $@

