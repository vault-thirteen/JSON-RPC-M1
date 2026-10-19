package jrm1

import (
	"encoding/json"
	"testing"

	"github.com/vault-thirteen/auxie/tester"
)

func Test_ParseParameters(t *testing.T) {
	aTest := tester.New(t)
	var re *RpcError
	var paramsRaw json.RawMessage
	initErrorMessages()

	// Test #1. Null parameters.
	{
		re = ParseParameters(nil, nil)
		aTest.MustBeDifferent(re, (*RpcError)(nil))
		aTest.MustBeEqual(re.Code, RpcErrorCode(-16))
	}

	// Test #2. Bad JSON.
	{
		paramsRaw = json.RawMessage([]byte(`x`))
		p := struct{}{}
		re = ParseParameters(&paramsRaw, &p)
		aTest.MustBeDifferent(re, (*RpcError)(nil))
		aTest.MustBeEqual(re.Code, RpcErrorCode(-16))
	}

	// Test #3. OK.
	{
		type P3 struct {
			N byte `json:"n"`
		}

		paramsRaw = json.RawMessage([]byte(`{"n":123}`))
		p := P3{}
		re = ParseParameters(&paramsRaw, &p)
		aTest.MustBeEqual(re, (*RpcError)(nil))
		aTest.MustBeEqual(p.N, byte(123))
	}

	// Test #4. Missing parameter.
	// This test is not possible in Go language while its JSON parser does not
	// check for missing fields, it just leaves them with their initial values,
	// i.e. pointers will be null if no field was found. Unfortunately, Go
	// language does not support strict parsing of JSON format.
	//{
	//	type P4 struct {
	//		Name    string `json:"name"`
	//		Surname string `json:"surname"`
	//	}
	//
	//	paramsRaw = json.RawMessage([]byte(`{"name":"John"}`))
	//	p := P4{}
	//	re = ParseParameters(&paramsRaw, &p)
	//	aTest.MustBeDifferent(re, (*RpcError)(nil))
	//}

	// Test #5. Unknown parameter.
	{
		type P5 struct {
			Name    string `json:"name"`
			Surname string `json:"surname"`
		}

		paramsRaw = json.RawMessage([]byte(`{"age":125}`))
		p := P5{}
		re = ParseParameters(&paramsRaw, &p)
		aTest.MustBeDifferent(re, (*RpcError)(nil))
	}
}
