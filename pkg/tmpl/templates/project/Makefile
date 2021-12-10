PKG=$(shell cat go.mod | grep "^module " | sed -e "s/module //g")
VERSION=v$(shell cat .version)

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOBUILD=CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-X ${PKG}/version.Version=${VERSION}"
OPENAPI=tools openapi

Workdir ?= ./cmd/{{ .Name }}

up: dockerize
	go run $(Workdir)

dockerize:
	go run $(Workdir) dockerize

migrate:
	go run $(Workdir) migrate

release:
	git push
	git push origin ${VERSION}
