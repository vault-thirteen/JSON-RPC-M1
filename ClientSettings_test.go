package jrm1

import (
	"testing"

	"github.com/vault-thirteen/auxie/tester"
)

func Test_NewClientSettings(t *testing.T) {
	aTest := tester.New(t)
	var cs *ClientSettings
	var err error

	// Test #1. Check error.
	cs, err = NewClientSettings("", "", 0, "", nil, nil, false)
	aTest.MustBeAnError(err)
	aTest.MustBeEqual(cs, (*ClientSettings)(nil))

	// Test #2. All clear.
	cs, err = NewClientSettings("http", "localhost", 80, "/", nil, nil, false)
	aTest.MustBeNoError(err)
	aTest.MustBeDifferent(cs, (*ClientSettings)(nil))
}

func Test_ClientSettings_Check(t *testing.T) {
	aTest := tester.New(t)
	var cs *ClientSettings
	var err error

	// Test #1. Schema is not set.
	cs = &ClientSettings{schema: "", host: "localhost", port: 80, path: "/"}
	err = cs.Check()
	aTest.MustBeAnError(err)

	// Test #2. Host is not set.
	cs = &ClientSettings{schema: "http", host: "", port: 80, path: "/"}
	err = cs.Check()
	aTest.MustBeAnError(err)

	// Test #3. Port is not set.
	cs = &ClientSettings{schema: "http", host: "localhost", port: 0, path: "/"}
	err = cs.Check()
	aTest.MustBeAnError(err)

	// Test #4. Path is not set.
	cs = &ClientSettings{schema: "http", host: "localhost", port: 80, path: ""}
	err = cs.Check()
	aTest.MustBeAnError(err)

	// Test #5. All clear.
	cs = &ClientSettings{schema: "http", host: "localhost", port: 80, path: "/"}
	err = cs.Check()
	aTest.MustBeNoError(err)
}
