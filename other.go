package jrm1

import (
	"bytes"
	"encoding/json"
)

// ParseParameters parses JSON bytes into an object representing parameters of
// called RPC method (function, procedure). This function must be called at the
// beginning of each RPC function to get parameters. Unfortunately, Go language
// can not do it automatically due to its technical limits.
func ParseParameters(params *json.RawMessage, dst any) (re *RpcError) {
	if params == nil {
		return NewRpcErrorFast(RpcErrorCode_InvalidParameters)
	}

	decoder := json.NewDecoder(bytes.NewReader(*params))
	decoder.DisallowUnknownFields()

	err := decoder.Decode(dst)
	if err != nil {
		return NewRpcErrorFast(RpcErrorCode_InvalidParameters)
	}

	return nil
}
