package jrm1

import (
	"testing"

	"github.com/vault-thirteen/auxie/tester"
)

func Test_ProcessorSettings_Check(t *testing.T) {
	aTest := tester.New(t)
	var ps *ProcessorSettings
	var err error

	// Test #1. Logging is enabled without catching.
	ps = &ProcessorSettings{
		LogExceptions:   true,
		CatchExceptions: false,
	}
	err = ps.Check()
	aTest.MustBeAnError(err)

	// Test #2. Duration field conflicts with ID field.
	someField := "abc"
	ps = &ProcessorSettings{
		DurationFieldName:  &someField,
		RequestIdFieldName: &someField,
	}
	err = ps.Check()
	aTest.MustBeAnError(err)

	// Test #3. All clear.
	someFieldA := "aa"
	someFieldB := "bb"
	ps = &ProcessorSettings{
		CatchExceptions:    true,
		LogExceptions:      true,
		CountRequests:      true,
		DurationFieldName:  &someFieldA,
		RequestIdFieldName: &someFieldB,
	}
	err = ps.Check()
	aTest.MustBeNoError(err)
}

func Test_ProcessorSettings_isDurationEnabled(t *testing.T) {
	aTest := tester.New(t)
	var ps *ProcessorSettings

	// Test #1.
	durField := "dur"
	ps = &ProcessorSettings{
		DurationFieldName: &durField,
	}
	aTest.MustBeEqual(ps.isDurationEnabled(), true)

	// Test #2.
	ps = &ProcessorSettings{
		DurationFieldName: nil,
	}
	aTest.MustBeEqual(ps.isDurationEnabled(), false)
}

func Test_ProcessorSettings_isRequestIdShown(t *testing.T) {
	aTest := tester.New(t)
	var ps *ProcessorSettings

	// Test #1.
	idField := "rid"
	ps = &ProcessorSettings{
		RequestIdFieldName: &idField,
	}
	aTest.MustBeEqual(ps.isRequestIdShown(), true)

	// Test #2.
	ps = &ProcessorSettings{
		RequestIdFieldName: nil,
	}
	aTest.MustBeEqual(ps.isRequestIdShown(), false)
}
