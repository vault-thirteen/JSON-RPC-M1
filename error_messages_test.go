package jrm1

import (
	"testing"

	"github.com/vault-thirteen/auxie/tester"
)

func Test_initErrorMessages(t *testing.T) {
	aTest := tester.New(t)

	// Test.
	errorMessages = nil
	initErrorMessages()
	aTest.MustBeEqual(len(errorMessages), 9)
}

func Test_findMessageForErrorCode(t *testing.T) {
	aTest := tester.New(t)
	var m RpcErrorMessage
	var err error

	initErrorMessages()

	// Test #1. Message is not found.
	m, err = findMessageForErrorCode(0)
	aTest.MustBeAnError(err)

	// Test #2. Message is found.
	m, err = findMessageForErrorCode(-32)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(m, RpcErrorMessage(RpcErrorMsg_InternalRpcError))
}
