package jrm1

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/vault-thirteen/auxie/errors"
)

const ProtocolNameM1 = "M1"

const (
	ErrFUnsupportedProtocol  = "unsupported protocol: %v"
	ErrRpcRequestIsMalformed = "RPC request is malformed"
)

// RpcRequest is a raw RPC request.
// Parameters in this raw request are not parsed.
type RpcRequest struct {
	// RPC protocol name.
	ProtocolName *string `json:"jsonrpc"`

	// Identifier of request.
	Id *string `json:"id"`

	// Name of the requested RPC function (method, procedure).
	Method *string `json:"method"`

	// Arguments for the requested RPC function (method, procedure).
	Parameters *json.RawMessage `json:"params"`
}

// NewRpcRequest is a constructor of a raw RPC request.
// It takes an input stream of bytes and decodes it using JSON format.
func NewRpcRequest(input io.ReadCloser) (rr *RpcRequest, err error) {
	defer func() {
		derr := input.Close()
		if derr != nil {
			rr = nil
			err = errors.Combine(err, derr)
		}
	}()

	rr = new(RpcRequest)
	err = json.NewDecoder(input).Decode(rr)
	if err != nil {
		return nil, err
	}

	return rr, nil
}

// HasAllRootFields tells if all the root fields are set, i.e. are not null
// pointers.
func (r *RpcRequest) HasAllRootFields() bool {
	if (r.ProtocolName == nil) ||
		(r.Id == nil) ||
		(r.Method == nil) ||
		(r.Parameters == nil) {
		return false
	}

	return true
}

// CheckProtocolVersion tells if the protocol version is correct.
func (r *RpcRequest) CheckProtocolVersion() (err error) {
	if r.ProtocolName == nil {
		return fmt.Errorf(ErrFUnsupportedProtocol, r.ProtocolName)
	}

	if *r.ProtocolName != ProtocolNameM1 {
		return fmt.Errorf(ErrFUnsupportedProtocol, *r.ProtocolName)
	}

	return nil
}
