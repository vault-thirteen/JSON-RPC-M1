package jrm1

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	mime "github.com/vault-thirteen/auxie/MIME"
	"github.com/vault-thirteen/auxie/header"
	"github.com/vault-thirteen/auxie/tester"
)

func Test_NewClient(t *testing.T) {
	aTest := tester.New(t)
	var cs *ClientSettings
	var c *Client
	var err error

	// Test #1. Error in settings.
	cs = &ClientSettings{}
	c, err = NewClient(cs)
	aTest.MustBeAnError(err)

	// Test #2. Normal settings.
	cs, err = NewClientSettings("http", "localhost", 80, "/", nil, nil, true)
	aTest.MustBeNoError(err)
	c, err = NewClient(cs)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(c.requestsCount.String(), "0")
}

func Test_Client_newHttpRequest(t *testing.T) {
	aTest := tester.New(t)
	var cs *ClientSettings
	var c *Client
	var rpcReq *RpcRequest
	var hr *http.Request
	var err error

	// Test.
	cs, err = NewClientSettings("http", "localhost", 80, "/", nil, nil, true)
	aTest.MustBeNoError(err)
	c, err = NewClient(cs)
	aTest.MustBeNoError(err)
	pn := ProtocolNameM1
	id := "123"
	m := "m"
	p := json.RawMessage([]byte(`{}`))
	rpcReq = &RpcRequest{
		ProtocolName: &pn,
		Id:           &id,
		Method:       &m,
		Parameters:   &p,
	}
	hr, err = c.newHttpRequest(context.Background(), rpcReq)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(hr.Host, "localhost:80")
	aTest.MustBeEqual(hr.URL.String(), "http://localhost:80/")
	aTest.MustBeEqual(hr.Method, http.MethodPost)
	aTest.MustBeEqual(hr.Header, http.Header{
		header.HttpHeaderContentType: []string{mime.TypeApplicationJson},
		header.HttpHeaderAccept:      []string{mime.TypeApplicationJson},
	})
	var reqBody []byte
	reqBody, err = io.ReadAll(hr.Body)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(string(reqBody), "{\n\t\"jsonrpc\": \"M1\",\n\t\"id\": \"123\",\n\t\"method\": \"m\",\n\t\"params\": {}\n}\n")
}

func Test_Client_call(t *testing.T) {
	aTest := tester.New(t)
	var cs *ClientSettings
	var c *Client
	var rpcReq *RpcRequest
	var err error

	cs, err = NewClientSettings("http", "localhost", 80, "/", nil, nil, true)
	aTest.MustBeNoError(err)
	c, err = NewClient(cs)
	aTest.MustBeNoError(err)

	pn := ProtocolNameM1
	pnx := "x"
	id := "123"
	m := "m"
	p := json.RawMessage([]byte(`{}`))

	// Test #1. 'HasAllRootFields' fails.
	rpcReq = &RpcRequest{
		ProtocolName: &pn,
		Id:           &id,
		Method:       &m,
		Parameters:   nil,
	}
	_, err = c.call(context.Background(), rpcReq)
	aTest.MustBeAnError(err)

	// Test #2. 'CheckProtocolVersion' fails.
	rpcReq = &RpcRequest{
		ProtocolName: &pnx,
		Id:           &id,
		Method:       &m,
		Parameters:   &p,
	}
	_, err = c.call(context.Background(), rpcReq)
	aTest.MustBeAnError(err)

	// Test #3. 'newHttpRequest' fails.
	// 'newHttpRequest' can not fail in practice, as it would require the
	// built-in JSON library to fail and other built-in Go libraries to fail.
	// We are not going to test the libraries which are built into the Go
	// language. This means, that we can not emulate this failure.

	// Test #4. 'httpClient.Do' fails.
	// RPC server must be shut down for this to succeed.
	rpcReq = &RpcRequest{
		ProtocolName: &pn,
		Id:           &id,
		Method:       &m,
		Parameters:   &p,
	}
	_, err = c.call(context.Background(), rpcReq)
	aTest.MustBeAnError(err)

	// Test #5. 'decoder.Decode' fails.
	// We will not emulate a broken server here while it will take forever to
	// implement.

	// Test #6. All clear.
	// We will not emulate a unit test here.
	// We do not waste time on deeds that are already done.
	// Use the unit test, located in the 'test' folder.
}

func Test_Client_GetRequestsCount(t *testing.T) {
	aTest := tester.New(t)
	var cs *ClientSettings
	var c *Client
	var err error

	// Test.
	cs, err = NewClientSettings("http", "localhost", 80, "/", nil, nil, true)
	aTest.MustBeNoError(err)
	c, err = NewClient(cs)
	aTest.MustBeNoError(err)
	c.incRequestsCount()
	aTest.MustBeEqual(c.GetRequestsCount(), "1")
}

func Test_Client_incRequestsCount(t *testing.T) {
	aTest := tester.New(t)
	var cs *ClientSettings
	var c *Client
	var err error

	// Test.
	cs, err = NewClientSettings("http", "localhost", 80, "/", nil, nil, true)
	aTest.MustBeNoError(err)
	c, err = NewClient(cs)
	aTest.MustBeNoError(err)
	c.incRequestsCount()
	aTest.MustBeEqual(c.requestsCount.String(), "1")
}
