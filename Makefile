
VERSION=$(shell cat VERSION)

build:
	CGO_ENABLED=0 go build -ldflags "-X github.com/drmdrew/syncrets/cmd.version=$(VERSION)"

