package jrm1

import (
	"testing"

	"github.com/vault-thirteen/auxie/tester"
)

func Test_RpcErrorMessage_Check(t *testing.T) {
	aTest := tester.New(t)
	var err error
	var rem RpcErrorMessage

	// Test #1. Positive.
	rem = RpcErrorMessage("ABBA")
	err = rem.Check()
	aTest.MustBeNoError(err)

	// Test #2. Negative.
	rem = RpcErrorMessage("")
	err = rem.Check()
	aTest.MustBeAnError(err)
}
