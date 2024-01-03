package main

import (
	"log"

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
		log.Fatal(err)
	}
}
