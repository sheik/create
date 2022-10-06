.PHONY: default install

default: package

Makefile: ;

install:
	@go install github.com/sheik/create/cmd/create@latest

%: install
	@create $@

