package jrm1

import (
	"fmt"
	"testing"

	"github.com/vault-thirteen/auxie/tester"
)

func Test_RpcErrorCode_Check(t *testing.T) {
	aTest := tester.New(t)
	var err error

	type TestDatum struct {
		errorCode       RpcErrorCode
		isErrorExpected bool
	}

	td := func(ec RpcErrorCode, iee bool) TestDatum {
		return TestDatum{errorCode: ec, isErrorExpected: iee}
	}

	testData := []TestDatum{
		// UGEC.
		// ...
		td(5, true),
		td(4, true),
		td(3, true),
		td(2, true),
		td(1, true),
		td(0, true),

		// RPC server errors.
		td(-1, false),
		td(-2, false),
		td(-4, false),
		td(-8, false),
		td(-16, false),
		td(-32, false),
		td(-64, false),
		td(-128, false),
		td(-256, false),

		// RPC server errors which are not implemented.
		td(-3, true),
		td(-5, true),
		td(-6, true),
		td(-7, true),
		td(-9, true),
		td(-10, true),
		td(-11, true),
		td(-12, true),
		td(-13, true),
		td(-14, true),
		td(-15, true),
		td(-17, true),
		// ...
		td(-63, true),
		td(-65, true),
		// ...
		td(-127, true),
		td(-129, true),
		// ...
		td(-255, true),
		td(-257, true),
		// ...

	}

	for i, testDatum := range testData {
		fmt.Printf("[%v]", i+1)

		err = testDatum.errorCode.Check()
		if testDatum.isErrorExpected {
			aTest.MustBeAnError(err)
		} else {
			aTest.MustBeNoError(err)
		}
	}
	fmt.Println()
}

func Test_RpcErrorCode_IsGeneratedByUser(t *testing.T) {
	aTest := tester.New(t)
	var result bool

	type TestDatum struct {
		errorCode      RpcErrorCode
		resultExpected bool
	}

	td := func(ec RpcErrorCode, re bool) TestDatum {
		return TestDatum{errorCode: ec, resultExpected: re}
	}

	testData := []TestDatum{
		td(-3, false),
		td(-2, false),
		td(-1, false),
		td(0, false),
		td(1, true),
		td(2, true),
		td(3, true),
	}

	for i, testDatum := range testData {
		fmt.Printf("[%v]", i+1)

		result = testDatum.errorCode.IsGeneratedByUser()
		aTest.MustBeEqual(result, testDatum.resultExpected)

	}
	fmt.Println()
}

func Test_RpcErrorCode_Int(t *testing.T) {
	aTest := tester.New(t)
	var rec RpcErrorCode

	// Test.
	rec = RpcErrorCode(123)
	aTest.MustBeEqual(rec.Int(), 123)
}
