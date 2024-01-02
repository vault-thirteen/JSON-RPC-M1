package jrm1

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	au "github.com/vault-thirteen/auxie/unicode"
)

const AsciiLowLine = '_'

const (
	ErrFBadSymbolInFunctionName = "bad symbol in function name: %v"
)

// RpcFunction represents a signature for an RPC function (method, procedure).
type RpcFunction func(params *json.RawMessage, metaData *ResponseMetaData) (result any, re *RpcError)

// GetName reads name of the RPC function (method, procedure).
//
// Do note that names of anonymous functions are not usable in practice while
// Go language has an unstable API and ABI.
//
// Go internal ABI specification is described at the following link:
// https://go.googlesource.com/go/+/refs/heads/dev.regabi/src/cmd/compile/internal-abi.md
//
// This document clearly states that ABI of Go language is unstable.
// -----------------------------------------------------------------------------
// This ABI is unstable and will change between Go versions. If you’re writing
// assembly code, please instead refer to Go’s assembly documentation, which
// describes Go’s stable ABI, known as ABI0.
// -----------------------------------------------------------------------------
func (f RpcFunction) GetName() string {
	fullName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	parts := strings.Split(fullName, ".")
	return parts[len(parts)-1]
}

// CheckFunctionName verifies name of a function.
func CheckFunctionName(fn string) (err error) {
	runes := []rune(fn)
	for _, r := range runes {
		if !isValidSymbolForFuncName(r) {
			return fmt.Errorf(ErrFBadSymbolInFunctionName, string(r))
		}
	}

	return nil
}

// isValidSymbolForFuncName verifies a symbol for a function name.
func isValidSymbolForFuncName(s rune) bool {
	return au.SymbolIsNumber(s) || au.SymbolIsLatLetter(s) || (s == AsciiLowLine)
}
