package main

import (
	"log"
	"syscall"

	"github.com/insidieux/pinchy/cmd/pinchy/internal"
	"github.com/sethvargo/go-signalcontext"
)

var (
	version string
)

func main() {
	ctx, cancel := signalcontext.On(syscall.SIGINT, syscall.SIGTERM)
	if cancel != nil {
		defer cancel()
	}
	if err := internal.NewCommand(version).ExecuteContext(ctx); err != nil {
		log.Fatalf(`Failed to execute command: %s`, err.Error())
	}
}
