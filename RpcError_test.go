package jrm1

import (
	"fmt"
	"testing"

	"github.com/vault-thirteen/auxie/tester"
)

func Test_NewRpcError(t *testing.T) {
	aTest := tester.New(t)
	var re, reExpected *RpcError
	var err error

	// Test #1. Invalid error code.
	initErrorMessages()
	re, err = NewRpcError(0, nil)
	aTest.MustBeAnError(err)
	aTest.MustBeEqual(err.Error(), "unsupported error code")

	// Test #2. Error code for user-generated errors.
	initErrorMessages()
	re, err = NewRpcError(1, nil)
	aTest.MustBeAnError(err)
	aTest.MustBeEqual(err.Error(), "user-generated errors use another constructor")

	// Test #3. Error message is not found.
	// It is not possible in practice, but we will emulate such an error.
	initErrorMessages()
	delete(errorMessages, -64)
	re, err = NewRpcError(-64, nil)
	aTest.MustBeAnError(err)
	aTest.MustBeEqual(err.Error(), "unknown error code: -64")

	// Test #4. Empty error message.
	// It is not possible in practice, but we will emulate such an error.
	initErrorMessages()
	errorMessages[RpcErrorCode(-64)] = RpcErrorMessage("")
	re, err = NewRpcError(-64, nil)
	aTest.MustBeAnError(err)
	aTest.MustBeEqual(err.Error(), "error message is not set")

	// Test #5. All clear.
	initErrorMessages()
	re, err = NewRpcError(-64, nil)
	aTest.MustBeNoError(err)
	reExpected = &RpcError{
		Code:    RpcErrorCode_ReservedForFuture_1,
		Message: RpcErrorMsg_ReservedForFuture_1,
	}
	aTest.MustBeEqual(re, reExpected)
}

func Test_NewRpcErrorFast(t *testing.T) {
	aTest := tester.New(t)
	var re, reExpected *RpcError
	var hasException bool

	// Test #1. No panic.
	initErrorMessages()
	re, hasException = _newRpcErrorFastWrapper(NewRpcErrorFast, -1)
	aTest.MustBeEqual(hasException, false)
	reExpected = &RpcError{
		Code:    RpcErrorCode_RequestIsNotReadable,
		Message: RpcErrorMsg_RequestIsNotReadable,
	}
	aTest.MustBeEqual(re, reExpected)

	// Test #2. Panic.
	initErrorMessages()
	re, hasException = _newRpcErrorFastWrapper(NewRpcErrorFast, 1)
	aTest.MustBeEqual(hasException, true)
	reExpected = (*RpcError)(nil)
	aTest.MustBeEqual(re, reExpected)
}

// _newRpcErrorFastWrapper runs the 'NewRpcErrorFast' and stops panic.
func _newRpcErrorFastWrapper(f func(code int) (re *RpcError), code int) (re *RpcError, hasException bool) {
	defer func() {
		x := recover()
		if x != nil {
			hasException = true
			fmt.Println(fmt.Sprintf("An exception was captured: %v", x))
		}
	}()

	re = f(code)

	return re, false
}

func Test_NewRpcErrorByUser(t *testing.T) {
	aTest := tester.New(t)
	var re, reExpected *RpcError
	var hasException bool

	// Test #1. No panic.
	initErrorMessages()
	re, hasException = _newRpcErrorByUserWrapper(NewRpcErrorByUser, 1, "Custom error message", 42)
	aTest.MustBeEqual(hasException, false)
	reExpected = &RpcError{
		Code:    1,
		Message: "Custom error message",
		Data:    42,
	}
	aTest.MustBeEqual(re, reExpected)

	// Test #2. Panic.
	initErrorMessages()
	re, hasException = _newRpcErrorByUserWrapper(NewRpcErrorByUser, -1, "x", nil)
	aTest.MustBeEqual(hasException, true)
	reExpected = (*RpcError)(nil)
	aTest.MustBeEqual(re, reExpected)
}

// _newRpcErrorByUserWrapper  runs the 'NewRpcErrorByUser' and stops panic.
func _newRpcErrorByUserWrapper(f func(code int, message string, data any) (re *RpcError), code int, message string, data any) (re *RpcError, hasException bool) {
	defer func() {
		x := recover()
		if x != nil {
			hasException = true
			fmt.Println(fmt.Sprintf("An exception was captured: %v", x))
		}
	}()

	re = f(code, message, data)

	return re, false
}
