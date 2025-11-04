#!/usr/bin/env gmake -f

GOOPTS=CGO_ENABLED=0
BUILDOPTS=-ldflags="-s -w" -a -gcflags=all=-l -trimpath

all: clean build

build:
	${GOOPTS} go build ${BUILDOPTS} -o aleesa-telegram-go ./cmd/aleesa-telegram-go

clean:
	rm -rf aleesa-telegram-go
	rm -rf settings-migrator

upgrade:
	rm -rf vendor
	go get -d -u -t ./...
	go mod tidy
	go mod vendor

migrator:
	${GOOPTS} go build ${BUILDOPTS} -o settings-migrator ./cmd/settings-migrator

# vim: set ft=make noet ai ts=4 sw=4 sts=4:
