package jrm1

import (
	"context"
	"encoding/json"
	"fmt"
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

// TODO
func Test_x(t *testing.T) {
	aTest := tester.New(t)
	var cs *ClientSettings
	var err error
	cs, err = NewClientSettings("http", "localhost", 80, "/", nil, nil, true)
	aTest.MustBeNoError(err)

	var c *Client
	c, err = NewClient(cs)
	aTest.MustBeNoError(err)

	var params = &SumParams{A: 245, B: 10}
	var result = new(SumResult)
	var rpcErr *RpcError
	ctx := context.Background()
	rpcErr, err = c.Call(ctx, "Sum", params, result)
	aTest.MustBeNoError(err)

	//TODO
	fmt.Println(result)
	fmt.Println(rpcErr)
}

func Test_newHttpRequest(t *testing.T) {
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

func Test_GetRequestsCount(t *testing.T) {
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

func Test_incRequestsCount(t *testing.T) {
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
