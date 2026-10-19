package main

import (
	"os"

	"github.com/vault-thirteen/auxie/errors"

	"github.com/vault-thirteen/JSON-RPC-M1/example/simple/a"
	"github.com/vault-thirteen/JSON-RPC-M1/example/simple/s"
)

func main() {
	app, err := a.NewApplication()
	mustBeNoError(err)

	app.Start()

	s.WaitForQuitSignalFromOS()

	err = app.Stop()
	mustBeNoError(err)
}

func mustBeNoError(err error) {
	if err != nil {
		errors.LogError(err)
		os.Exit(1)
	}
}
