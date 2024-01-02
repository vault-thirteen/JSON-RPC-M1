package jrm1

import (
	"fmt"
	"testing"

	"github.com/vault-thirteen/auxie/tester"
)

func Test_RpcFunction_GetName(t *testing.T) {
	aTest := tester.New(t)
	rpcFn := RpcFunction(RpcFunctionExampleOne)
	fname := rpcFn.GetName()
	aTest.MustBeEqual(fname, "RpcFunctionExampleOne")
}

func Test_CheckFunctionName(t *testing.T) {
	aTest := tester.New(t)
	aTest.MustBeNoError(CheckFunctionName("0123456789"))
	aTest.MustBeNoError(CheckFunctionName("abcdefghijklmnopqrstuvwxyz"))
	aTest.MustBeNoError(CheckFunctionName("ABCDEFGHIJKLMNOPQRSTUVWXYZ"))
	aTest.MustBeNoError(CheckFunctionName("_"))
	aTest.MustBeAnError(CheckFunctionName("*"))
}

func Test_isValidSymbolForFuncName(t *testing.T) {
	aTest := tester.New(t)

	validSymbols := []rune{}

	for _, s := range []rune("abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "0123456789" + "_") {
		validSymbols = append(validSymbols, s)
	}

	aTest.MustBeEqual(len(validSymbols), 26+26+10+1)

	for _, s := range validSymbols {
		fmt.Printf("%s", string(s))
		aTest.MustBeEqual(isValidSymbolForFuncName(s), true)
	}
	fmt.Println()

	aTest.MustBeEqual(isValidSymbolForFuncName('*'), false)
	aTest.MustBeEqual(isValidSymbolForFuncName('.'), false)
	aTest.MustBeEqual(isValidSymbolForFuncName(' '), false)
}
