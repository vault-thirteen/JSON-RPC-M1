package s

import (
	"encoding/json"
	"fmt"
	"time"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
)

// SumParams are parameters for the 'Sum' function.
type SumParams struct {
	A byte `json:"a"`
	B byte `json:"b"`
}

// SumResult are results of the 'Sum' function.
type SumResult struct {
	C byte `json:"c"`
}

// Sum simply sums two bytes.
func Sum(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p SumParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	// We are catching the possible overflow here.
	if (255 - p.A) < p.B {
		return nil, jrm1.NewRpcErrorByUser(1, "overflow", p)
	}

	r := new(SumResult)
	r.C = p.A + p.B

	return r, nil
}

// Crash shows the request ID, adds a meta-data field, emulates a long process
// and crashes. This function shows how internal RPC errors look like.
func Crash(_ *json.RawMessage, meta *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	requestId := meta.GetField(RequestIdFieldName)
	fmt.Println(fmt.Sprintf("In called function, RID=%v.", requestId))
	meta.AddFieldFast("Laboratory", "L3")
	time.Sleep(time.Second * 1)

	x := 1
	x = x / (x - x)

	return x, nil
}
