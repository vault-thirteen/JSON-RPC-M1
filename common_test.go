package jrm1

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type _badCloser struct {
	io.Reader
}

var _metaFieldName_RID = "rid"
var _metaFieldName_Duration = "dur"

func (bc _badCloser) Close() error { return errors.New("close error") }

func RpcFunctionExampleOne(params *json.RawMessage, metaData *ResponseMetaData) (result any, re *RpcError) {
	fmt.Println(fmt.Sprintf("Raw parameters: %v, Meta-data: %v.", params, metaData))
	return nil, nil
}

func RpcFunctionExampleTwo(_ *json.RawMessage, metaData *ResponseMetaData) (result any, re *RpcError) {
	// Here we are trying to break the meta-data set by the framework.
	err := metaData.RemoveField(_metaFieldName_RID)
	if err != nil {
		fmt.Println(err)
	}
	return nil, nil
}

func RpcFunctionExampleThree(_ *json.RawMessage, metaData *ResponseMetaData) (result any, re *RpcError) {
	// Here we are trying to break the meta-data set by the framework.
	err := metaData.AddField(_metaFieldName_Duration, "x")
	if err != nil {
		fmt.Println(err)
	}
	return nil, nil
}

func RpcFunctionExample四(_ *json.RawMessage, _ *ResponseMetaData) (result any, re *RpcError) {
	fmt.Println(`こんにちは`)
	return nil, nil
}

func RpcFunctionExampleFive(_ *json.RawMessage, _ *ResponseMetaData) (result any, re *RpcError) {
	return 2024, nil
}

func RpcFunctionExampleCrasher(_ *json.RawMessage, _ *ResponseMetaData) (result any, re *RpcError) {
	x := 1
	x = x / (x - x)
	return `haha`, nil
}

// SumParams are parameters for the 'Sum' function.
type SumParams struct {
	A byte `json:"a"`
	B byte `json:"b"`
}

// SumResult are results of the 'Sum' function.
type SumResult struct {
	C byte `json:"c"`
}

// RpcFunctionSum simply sums two bytes.
func RpcFunctionSum(params *json.RawMessage, _ *ResponseMetaData) (result any, re *RpcError) {
	var p SumParams
	re = ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	// We are catching the possible overflow here.
	if (255 - p.A) < p.B {
		return nil, NewRpcErrorByUser(1, "overflow", p)
	}

	r := new(SumResult)
	r.C = p.A + p.B

	return r, nil
}
