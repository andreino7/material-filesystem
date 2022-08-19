# This is a Makefile "target".  The first target is always the default one.
help:
	@echo "Welcome to material-filesystem!  Here's a list of available Makefile targets:"
	@$(MAKE) list-targets

# Lists all available targets within the Makefile, per https://stackoverflow.com/a/26339924
.PHONY: list-targets
list-targets:
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'

build-fs:
	go build -o ./build/fs ./fs/cmd/main.go

build-daemon:
	go build -o ./build/fs-daemon ./daemon/cmd/main.go

test:
	go test ./... -cover

vet:
	go vet ./...

staticcheck:
	staticcheck ./...
