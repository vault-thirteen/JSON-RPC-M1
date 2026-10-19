package jrm1

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/vault-thirteen/auxie/tester"
)

func Test_NewRpcRequest(t *testing.T) {
	aTest := tester.New(t)
	var err error
	var rr, rrExpected *RpcRequest
	var data io.ReadCloser

	// Test #1. Positive.
	dataStr := `{"id": "12345"}`
	data = io.NopCloser(bytes.NewReader([]byte(dataStr)))
	id := "12345"
	rrExpected = &RpcRequest{Id: &id}
	//
	rr, err = NewRpcRequest(data)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(rr, rrExpected)

	// Test #2. Negative. Not a JSON format.
	dataStr = `x`
	data = io.NopCloser(bytes.NewReader([]byte(dataStr)))
	rrExpected = (*RpcRequest)(nil)
	//
	rr, err = NewRpcRequest(data)
	aTest.MustBeAnError(err)
	aTest.MustBeEqual(rr, rrExpected)

	// Test #3. Negative. Bad closing.
	dataStr = `{}`
	data = _badCloser{bytes.NewBufferString(dataStr)}
	rrExpected = (*RpcRequest)(nil)
	//
	rr, err = NewRpcRequest(data)
	aTest.MustBeAnError(err)
	aTest.MustBeEqual(rr, rrExpected)
}

func Test_RpcRequest_HasAllRootFields(t *testing.T) {
	aTest := tester.New(t)
	var rr *RpcRequest

	// Common data.
	pn := "PN"
	id := "12345"
	m := "m"
	p := json.RawMessage([]byte{})

	// Test #1. Positive.
	rr = &RpcRequest{ProtocolName: &pn, Id: &id, Method: &m, Parameters: &p}
	aTest.MustBeEqual(rr.HasAllRootFields(), true)

	// Test #2. Field 1 is null.
	rr = &RpcRequest{ProtocolName: nil, Id: &id, Method: &m, Parameters: &p}
	aTest.MustBeEqual(rr.HasAllRootFields(), false)

	// Test #3. Field 2 is null.
	rr = &RpcRequest{ProtocolName: &pn, Id: nil, Method: &m, Parameters: &p}
	aTest.MustBeEqual(rr.HasAllRootFields(), false)

	// Test #4. Field 3 is null.
	rr = &RpcRequest{ProtocolName: &pn, Id: &id, Method: nil, Parameters: &p}
	aTest.MustBeEqual(rr.HasAllRootFields(), false)

	// Test #5. Field 4 is null.
	rr = &RpcRequest{ProtocolName: &pn, Id: &id, Method: &m, Parameters: nil}
	aTest.MustBeEqual(rr.HasAllRootFields(), false)
}

func Test_RpcRequest_CheckProtocolVersion(t *testing.T) {
	aTest := tester.New(t)
	var rr *RpcRequest
	var err error

	// Common data.
	pn := ProtocolNameM1

	// Test #1. Positive.
	rr = &RpcRequest{ProtocolName: &pn}
	err = rr.CheckProtocolVersion()
	aTest.MustBeNoError(err)

	// Test #2. Protocol is not set.
	rr = &RpcRequest{}
	err = rr.CheckProtocolVersion()
	aTest.MustBeAnError(err)

	// Test #3. Protocol is not correct.
	pn2 := "x"
	rr = &RpcRequest{ProtocolName: &pn2}
	err = rr.CheckProtocolVersion()
	aTest.MustBeAnError(err)
}
