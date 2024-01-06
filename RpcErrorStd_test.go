package jrm1

import (
	"testing"

	"github.com/vault-thirteen/auxie/tester"
)

func Test_RpcErrorStd_Error(t *testing.T) {
	aTest := tester.New(t)
	var res *RpcErrorStd
	var err error

	// Test.
	res = &RpcErrorStd{
		RpcError{
			Code:    1,
			Message: "One",
			Data:    true,
		},
	}
	aTest.MustBeEqual(res.Error(), "One")

	err = res
	aTest.MustBeAnError(err)
	aTest.MustBeEqual(err.Error(), "One")

	var ok bool
	res, ok = err.(*RpcErrorStd)
	aTest.MustBeEqual(ok, true)
	aTest.MustBeEqual(res.Code, RpcErrorCode(1))
}
