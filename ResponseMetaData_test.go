package jrm1

import (
	"fmt"
	"testing"

	"github.com/vault-thirteen/auxie/tester"
)

func Test_ResponseMetaData_AddField(t *testing.T) {
	aTest := tester.New(t)
	var md ResponseMetaData
	var err error

	// Test #1. Normal addition.
	md = make(ResponseMetaData)
	err = md.AddField("key", "value")
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(md["key"], "value")

	// Test #2. Duplicate key.
	md = make(ResponseMetaData)
	err = md.AddField("key", "value-1")
	aTest.MustBeNoError(err)
	err = md.AddField("key", "value-2")
	aTest.MustBeAnError(err)
}

// _newRpcErrorFastWrapper runs the 'NewRpcErrorFast' and stops panic.
func _addFieldFastWrapper(f func(key string, value any), key string, value any) (hasException bool) {
	defer func() {
		x := recover()
		if x != nil {
			hasException = true
			fmt.Println(fmt.Sprintf("An exception was captured: %v", x))
		}
	}()

	f(key, value)

	return false
}

func Test_ResponseMetaData_AddFieldFast(t *testing.T) {
	aTest := tester.New(t)
	var hasException bool
	var md ResponseMetaData

	// Test #1. No panic.
	md = make(ResponseMetaData)
	hasException = _addFieldFastWrapper(md.AddFieldFast, "a", 123)
	aTest.MustBeEqual(hasException, false)

	// Test #2. Panic.
	md = make(ResponseMetaData)
	hasException = _addFieldFastWrapper(md.AddFieldFast, "a", 123)
	aTest.MustBeEqual(hasException, false)
	hasException = _addFieldFastWrapper(md.AddFieldFast, "a", 456)
	aTest.MustBeEqual(hasException, true)
}

func Test_ResponseMetaData_GetField(t *testing.T) {
	aTest := tester.New(t)
	var md ResponseMetaData
	var err error

	// Test.
	md = make(ResponseMetaData)
	err = md.AddField("key", "value")
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(md.GetField("key"), "value")
}

func Test_ResponseMetaData_RemoveField(t *testing.T) {
	aTest := tester.New(t)
	var md ResponseMetaData
	var err error

	// Test #1. Normal deletion.
	md = make(ResponseMetaData)
	err = md.AddField("key", "value")
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(md["key"], "value")
	err = md.RemoveField("key")
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(md["key"], nil)

	// Test #2. Absent key.
	md = make(ResponseMetaData)
	err = md.AddField("key", "value")
	aTest.MustBeNoError(err)
	err = md.RemoveField("key")
	aTest.MustBeNoError(err)
	err = md.RemoveField("key")
	aTest.MustBeAnError(err)
}
