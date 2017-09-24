
VERSION=$(shell cat VERSION)

build: syncrets
	go build -ldflags "-X github.com/drmdrew/syncrets/cmd.version=$(VERSION)"

