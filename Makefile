.PHONY: default install

default: package

%:
	@go install ./cmd/create
	@create $@
