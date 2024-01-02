package jrm1

import "errors"

// RPC error messages.
const (
	RpcErrorMsg_RequestIsNotReadable = "Request is not readable"
	RpcErrorMsg_InvalidRequest       = "Invalid request"
	RpcErrorMsg_UnsupportedProtocol  = "Unsupported protocol"
	RpcErrorMsg_UnknownMethod        = "Unknown method"
	RpcErrorMsg_InvalidParameters    = "Invalid parameters"
	RpcErrorMsg_InternalRpcError     = "Internal RPC error"
	//
	RpcErrorMsg_ReservedForFuture_1 = "Reserved for future (1)"
	RpcErrorMsg_ReservedForFuture_2 = "Reserved for future (2)"
	RpcErrorMsg_ReservedForFuture_3 = "Reserved for future (3)"

	RpcErrorMsg_Empty = ""
)

const (
	ErrErrorMessageIsNotSet = "error message is not set"
)

// RpcErrorMessage is a message for an RPC error.
// It stores some textual description of an error.
type RpcErrorMessage string

// Check performs a basic check – it ensures that message is set, i.e. is not
// empty. This function does not validate correctness of the text.
func (rem RpcErrorMessage) Check() (err error) {
	if len(rem) == 0 {
		return errors.New(ErrErrorMessageIsNotSet)
	}

	return nil
}
