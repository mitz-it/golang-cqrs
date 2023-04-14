#!make

.PHONY: test

test:
	mkdir -p out && \
	go test -covermode=atomic -coverprofile=out/coverage.out $(go list ./... | grep -v vendor/)