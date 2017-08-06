

build: syncrets
	go build -ldflags "-X github.com/drmdrew/syncrets/cmd.version=0.0.1-1"

