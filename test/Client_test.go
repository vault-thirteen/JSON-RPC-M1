package test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/vault-thirteen/JSON-RPC-M1"
	"github.com/vault-thirteen/JSON-RPC-M1/example/simple/a"
	"github.com/vault-thirteen/JSON-RPC-M1/example/simple/s"
	"github.com/vault-thirteen/auxie/tester"
)

// This file tests two public functions of the client: Call and CallRaw.
// It uses the provided simple example as an RPC server to emulate responses.
// This test can be considered as a simple automated unit test.
// For manual unit tests, see a description of the simple example.

func Test_Client_Call(t *testing.T) {
	aTest := tester.New(t)

	srvApp, err := a.NewApplication()
	aTest.MustBeNoError(err)
	srvApp.Start()

	defer func() {
		derr := srvApp.Stop()
		aTest.MustBeNoError(derr)
	}()

	// Test.
	var cs *jrm1.ClientSettings
	cs, err = jrm1.NewClientSettings("http", "localhost", 80, "/", nil, nil, true)
	aTest.MustBeNoError(err)

	var c *jrm1.Client
	c, err = jrm1.NewClient(cs)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(c.GetRequestsCount(), "0")

	p := &s.SumParams{A: 1, B: 2}
	res := new(s.SumResult)

	var re *jrm1.RpcError
	re, err = c.Call(context.Background(), "Sum", p, res)
	aTest.MustBeNoError(err)

	aTest.MustBeEqual(re, (*jrm1.RpcError)(nil))
	aTest.MustBeEqual(res, &s.SumResult{C: 3})
	aTest.MustBeEqual(c.GetRequestsCount(), "1")
}

func Test_Client_CallRaw(t *testing.T) {
	aTest := tester.New(t)

	srvApp, err := a.NewApplication()
	aTest.MustBeNoError(err)
	srvApp.Start()

	defer func() {
		derr := srvApp.Stop()
		aTest.MustBeNoError(derr)
	}()

	// Test.
	var cs *jrm1.ClientSettings
	cs, err = jrm1.NewClientSettings("http", "localhost", 80, "/", nil, nil, true)
	aTest.MustBeNoError(err)

	var c *jrm1.Client
	c, err = jrm1.NewClient(cs)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(c.GetRequestsCount(), "0")

	pn := jrm1.ProtocolNameM1
	id := "123"
	m := "Sum"
	p := json.RawMessage([]byte(`{"a": 1, "b": 2}`))
	rpcReq := &jrm1.RpcRequest{
		ProtocolName: &pn,
		Id:           &id,
		Method:       &m,
		Parameters:   &p,
	}

	var rpcResp *jrm1.RpcResponseRaw
	rpcResp, err = c.CallRaw(context.Background(), rpcReq)
	aTest.MustBeNoError(err)

	aTest.MustBeEqual(rpcResp.ProtocolName, jrm1.ProtocolNameM1)
	aTest.MustBeEqual(*rpcResp.Id, "123")
	aTest.MustBeEqual(string(*rpcResp.Result), `{"c":3}`)
	aTest.MustBeEqual(rpcResp.Error, (*jrm1.RpcError)(nil))
	aTest.MustBeEqual(rpcResp.Meta, &jrm1.ResponseMetaData{"dur": json.Number("0")})
	aTest.MustBeEqual(rpcResp.OK, true)
	aTest.MustBeEqual(c.GetRequestsCount(), "1")
}
