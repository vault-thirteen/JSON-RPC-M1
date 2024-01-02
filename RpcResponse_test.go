package jrm1

import (
	"testing"

	"github.com/vault-thirteen/auxie/tester"
)

func Test_NewRpcResponse(t *testing.T) {
	aTest := tester.New(t)
	var resp, respExpected *RpcResponse

	// Test.
	md := make(ResponseMetaData)
	respExpected = &RpcResponse{
		ProtocolName: ProtocolNameM1,
		Id:           nil,
		Result:       nil,
		Error:        nil,
		Meta:         &md,
		OK:           false,
	}
	resp = NewRpcResponse()
	aTest.MustBeEqual(resp, respExpected)
}

func Test_RpcResponse_hasError(t *testing.T) {
	aTest := tester.New(t)
	var resp *RpcResponse

	// Test #1. Error is not set.
	resp = &RpcResponse{Error: nil}
	aTest.MustBeEqual(resp.hasError(), false)

	// Test #2. Error is set.
	re := RpcError{}
	resp = &RpcResponse{Error: &re}
	aTest.MustBeEqual(resp.hasError(), true)
}
