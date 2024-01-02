package jrm1

import "encoding/json"

// RpcResponseRaw is a raw RPC response.
type RpcResponseRaw struct {
	// RPC protocol name.
	ProtocolName string `json:"jsonrpc"`

	// Identifier of request and its response.
	Id *string `json:"id"`

	// Result returned by the called RPC function (method, procedure).
	Result *json.RawMessage `json:"result"`

	// Error returned by the called RPC function (method, procedure).
	Error *RpcError `json:"error"`

	// Additional meta-information about the RPC call. This information is not
	// returned by the called function (method, procedure), however the called
	// function (method, procedure) has an interface to interact with this
	// information. It is provided for some specific purposes, such as
	// debugging, logging, transmission of various environment settings, not
	// directly related to the RPC function (method, procedure). The most
	// evident usage of this object is providing an ability to see an ID of the
	// RPC request inside the called RPC function (method, procedure). Another
	// usage example is storing duration of the RPC call, which can be
	// automatically performed by this framework.
	Meta *ResponseMetaData `json:"meta,omitempty"`

	// Flag of success.
	// This field serves at least two purposes. First, it provides additional
	// robustness via redundancy. Second, it allows to check for success not
	// having to parse the whole result JSON object. If you are familiar with
	// the HTTP status code 200, you will understand the idea.
	OK bool `json:"ok"`
}
