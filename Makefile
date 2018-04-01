
VERSION=$(shell cat VERSION)

build:
	CGO_ENABLED=0 go build -a -ldflags "-X github.com/drmdrew/syncrets/cmd.version=$(VERSION)"

