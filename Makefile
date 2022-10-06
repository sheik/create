.PHONY: default install

default:
	@go install github.com/sheik/create/cmd/create@latest
	@$(MAKE) -s package

%:
	@create $@
