package jrm1

import (
	"context"
	"fmt"
	"testing"

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
