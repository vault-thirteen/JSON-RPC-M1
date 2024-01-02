package jrm1

import "errors"

const (
	ErrUserGeneratedErrorsUseAnotherConstructor = "user-generated errors use another constructor"
	ErrUserGeneratedErrorsHaveSpecialCodes      = "user-generated errors have special codes"
)

// RpcError is an RPC error.
type RpcError struct {
	// Error code.
	Code RpcErrorCode `json:"code"`

	// Error message.
	Message RpcErrorMessage `json:"message"`

	// Some additional details about an error.
	Data any `json:"data"`
}

// NewRpcError is a common constructor of an RPC error.
// It returns an error in a Go-language style.
// This constructor is not for user-generated errors.
func NewRpcError(code int, data any) (re *RpcError, err error) {
	re = &RpcError{
		Code: RpcErrorCode(code),
	}

	err = re.Code.Check()
	if err != nil {
		if re.Code.IsGeneratedByUser() {
			return nil, errors.New(ErrUserGeneratedErrorsUseAnotherConstructor)
		}

		return nil, err
	}

	re.Message, err = findMessageForErrorCode(re.Code)
	if err != nil {
		return nil, err
	}

	err = re.Message.Check()
	if err != nil {
		return nil, err
	}

	re.Data = data

	return re, nil
}

// NewRpcErrorFast is a fast constructor which throws an exception on error.
// This constructor should only be used if you know what it is.
func NewRpcErrorFast(code int) (re *RpcError) {
	var err error
	re, err = NewRpcError(code, nil)
	if err != nil {
		panic(err)
		return nil // This line is normally unreachable, but we expect anything.
	}

	return re
}

// NewRpcErrorByUser is a constructor of an RPC error created by user.
// This function throws an exception on error.
func NewRpcErrorByUser(code int, message string, data any) (re *RpcError) {
	re = &RpcError{
		Code:    RpcErrorCode(code),
		Message: RpcErrorMessage(message),
		Data:    data,
	}

	if !re.Code.IsGeneratedByUser() {
		err := errors.New(ErrUserGeneratedErrorsHaveSpecialCodes)
		panic(err)
		return nil // This line is normally unreachable, but we expect anything.
	}

	return re
}
