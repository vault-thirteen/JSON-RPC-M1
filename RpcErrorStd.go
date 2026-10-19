package jrm1

// RpcErrorStd is a clone of an RPC error which supports methods of the
// standard error interface.
type RpcErrorStd struct {
	RpcError
}

func NewRpcErrorStd(code RpcErrorCode, msg RpcErrorMessage, data any) (res *RpcErrorStd) {
	return &RpcErrorStd{
		RpcError{
			Code:    code,
			Message: msg,
			Data:    data,
		},
	}
}

// Error is a standard method of the error interface.
func (res *RpcErrorStd) Error() string {
	return res.Message.String()
}
