package jrm1

import "fmt"

// errorMessages is a mapping of built-in RPC error messages by their error
// code.
var errorMessages map[RpcErrorCode]RpcErrorMessage

// initErrorMessages creates a mapping of built-in RPC error messages by their
// error code.
func initErrorMessages() {
	errorMessages = map[RpcErrorCode]RpcErrorMessage{
		RpcErrorCode_RequestIsNotReadable: RpcErrorMsg_RequestIsNotReadable,
		RpcErrorCode_InvalidRequest:       RpcErrorMsg_InvalidRequest,
		RpcErrorCode_UnsupportedProtocol:  RpcErrorMsg_UnsupportedProtocol,
		RpcErrorCode_UnknownMethod:        RpcErrorMsg_UnknownMethod,
		RpcErrorCode_InvalidParameters:    RpcErrorMsg_InvalidParameters,
		RpcErrorCode_InternalRpcError:     RpcErrorMsg_InternalRpcError,
		RpcErrorCode_ReservedForFuture_1:  RpcErrorMsg_ReservedForFuture_1,
		RpcErrorCode_ReservedForFuture_2:  RpcErrorMsg_ReservedForFuture_2,
		RpcErrorCode_ReservedForFuture_3:  RpcErrorMsg_ReservedForFuture_3,
	}
}

// findMessageForErrorCode gets an error message from the map of built-in
// error messages by its error code.
func findMessageForErrorCode(code RpcErrorCode) (message RpcErrorMessage, err error) {
	var ok bool
	message, ok = errorMessages[code]
	if !ok {
		return RpcErrorMsg_Empty, fmt.Errorf(ErrFUnknownErrorCode, code)
	}

	return message, nil
}
