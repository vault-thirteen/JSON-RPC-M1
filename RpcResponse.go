package jrm1

// RpcResponse is an RPC response.
type RpcResponse struct {
	// RPC protocol name.
	ProtocolName string `json:"jsonrpc"`

	// Identifier of request and its response.
	Id *string `json:"id"`

	// Result returned by the called RPC function (method, procedure).
	Result any `json:"result"`

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

// NewRpcResponse creates an empty (template) RPC response. Do note that the
// returned response is not valid until it is filled with actual data.
func NewRpcResponse() (resp *RpcResponse) {
	md := make(ResponseMetaData)

	return &RpcResponse{
		ProtocolName: ProtocolNameM1,
		Id:           nil,
		Result:       nil,
		Error:        nil,
		Meta:         &md,
		OK:           false,
	}
}

// hasError tells if the response finished with an error or not.
func (r *RpcResponse) hasError() bool {
	return r.Error != nil
}
