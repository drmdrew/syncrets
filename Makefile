
VERSION=$(shell cat VERSION)

build:
	go build -ldflags "-X github.com/drmdrew/syncrets/cmd.version=$(VERSION)"

