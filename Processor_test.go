package jrm1

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"testing"

	mime "github.com/vault-thirteen/auxie/MIME"
	"github.com/vault-thirteen/auxie/header"
	hh "github.com/vault-thirteen/auxie/http-helper"
	"github.com/vault-thirteen/auxie/tester"
)

func Test_NewProcessor(t *testing.T) {
	aTest := tester.New(t)
	var ps *ProcessorSettings
	var p *Processor
	var err error

	// Test #1. Bad settings.
	ps = &ProcessorSettings{
		LogExceptions:   true,
		CatchExceptions: false,
	}
	p, err = NewProcessor(ps)
	aTest.MustBeAnError(err)
	aTest.MustBeEqual(p, (*Processor)(nil))

	// Test #2. Normal settings. Error messages are not initialised.
	ps = &ProcessorSettings{
		LogExceptions:   true,
		CatchExceptions: true,
	}
	errorMessages = nil
	aTest.MustBeEqual(len(errorMessages), 0)
	p, err = NewProcessor(ps)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(p.settings, ps)
	aTest.MustBeEqual(len(errorMessages), 9)
}

func Test_Processor_AddFunc(t *testing.T) {
	aTest := tester.New(t)
	var ps *ProcessorSettings
	var p *Processor
	var err error

	// Test #1. Function name is not valid.
	ps = &ProcessorSettings{}
	p, err = NewProcessor(ps)
	aTest.MustBeNoError(err)
	err = p.AddFunc(RpcFunctionExample四)
	aTest.MustBeAnError(err)
	aTest.MustBeEqual(err.Error(), `bad symbol in function name: 四`)

	// Test #2. Normal function and duplicate function.
	ps = &ProcessorSettings{}
	p, err = NewProcessor(ps)
	aTest.MustBeNoError(err)
	err = p.AddFunc(RpcFunctionExampleOne)
	aTest.MustBeNoError(err)
	err = p.AddFunc(RpcFunctionExampleOne)
	aTest.MustBeAnError(err)
	aTest.MustBeEqual(err.Error(), `duplicate function`)
}

func Test_Processor_AddFuncFast(t *testing.T) {
	aTest := tester.New(t)
	var ps *ProcessorSettings
	var p *Processor
	var err error
	var hasException bool

	// Test #1. All clear.
	ps = &ProcessorSettings{}
	p, err = NewProcessor(ps)
	aTest.MustBeNoError(err)
	hasException = _Processor_AddFuncFast_Wrapper(p.AddFuncFast, RpcFunctionExampleOne)
	aTest.MustBeEqual(hasException, false)

	// Test #2. Exception.
	ps = &ProcessorSettings{}
	p, err = NewProcessor(ps)
	aTest.MustBeNoError(err)
	hasException = _Processor_AddFuncFast_Wrapper(p.AddFuncFast, RpcFunctionExample四)
	aTest.MustBeEqual(hasException, true)
}

// _Processor_AddFuncFast_Wrapper  runs the 'AddFuncFast' and stops panic.
func _Processor_AddFuncFast_Wrapper(testedFn func(RpcFunction), arg RpcFunction) (hasException bool) {
	defer func() {
		x := recover()
		if x != nil {
			hasException = true
			fmt.Println(fmt.Sprintf("An exception was captured: %v", x))
		}
	}()

	testedFn(arg)

	return false
}

func Test_Processor_RemoveFunc(t *testing.T) {
	aTest := tester.New(t)
	var ps *ProcessorSettings
	var p *Processor
	var err error

	// Test #1. Function is found.
	ps = &ProcessorSettings{}
	p, err = NewProcessor(ps)
	aTest.MustBeNoError(err)
	err = p.AddFunc(RpcFunctionExampleOne)
	aTest.MustBeNoError(err)
	err = p.RemoveFunc("RpcFunctionExampleOne")
	aTest.MustBeNoError(err)

	// Test #2. Function is not found.
	ps = &ProcessorSettings{}
	p, err = NewProcessor(ps)
	aTest.MustBeNoError(err)
	err = p.RemoveFunc("RpcFunctionExampleOne")
	aTest.MustBeAnError(err)
	aTest.MustBeEqual(err.Error(), `function is not found`)
}

func Test_Processor_FindFunc(t *testing.T) {
	aTest := tester.New(t)
	var ps *ProcessorSettings
	var p *Processor
	var err error

	// Test #1. Function is found.
	ps = &ProcessorSettings{}
	p, err = NewProcessor(ps)
	aTest.MustBeNoError(err)
	err = p.AddFunc(RpcFunctionExampleOne)
	aTest.MustBeNoError(err)
	err = p.FindFunc("RpcFunctionExampleOne")
	aTest.MustBeNoError(err)

	// Test #2. Function is not found.
	ps = &ProcessorSettings{}
	p, err = NewProcessor(ps)
	aTest.MustBeNoError(err)
	err = p.FindFunc("RpcFunctionExampleOne")
	aTest.MustBeAnError(err)
	aTest.MustBeEqual(err.Error(), `function is not found`)
}

func Test_Processor_RunFunc(t *testing.T) {
	aTest := tester.New(t)
	var ps *ProcessorSettings
	var p *Processor
	var err error
	var result any
	var re *RpcError

	// Test #1. Function is not found.
	{
		ps = &ProcessorSettings{}
		p, err = NewProcessor(ps)
		aTest.MustBeNoError(err)
		result, re = p.RunFunc("RpcFunctionExampleOne", nil, nil)
		aTest.MustBeDifferent(re, (*RpcError)(nil))
		aTest.MustBeEqual(re.Code, RpcErrorCode(-8))
		aTest.MustBeEqual(result, nil)
	}

	// Test #2. All clear.
	{
		ps = &ProcessorSettings{}
		p, err = NewProcessor(ps)
		aTest.MustBeNoError(err)
		err = p.AddFunc(RpcFunctionExampleFive)
		aTest.MustBeNoError(err)
		result, re = p.RunFunc("RpcFunctionExampleFive", nil, nil)
		aTest.MustBeEqual(result, 2024)
		aTest.MustBeEqual(re, (*RpcError)(nil))
	}

	// Test #3. Exception.
	{
		ps = &ProcessorSettings{
			CatchExceptions: true,
			LogExceptions:   true,
		}
		p, err = NewProcessor(ps)
		aTest.MustBeNoError(err)
		err = p.AddFunc(RpcFunctionExampleCrasher)
		aTest.MustBeNoError(err)
		result, re = p.RunFunc("RpcFunctionExampleCrasher", nil, nil)
		aTest.MustBeEqual(result, nil)
		aTest.MustBeDifferent(re, (*RpcError)(nil))
		aTest.MustBeEqual(re.Code, RpcErrorCode(-32))
	}
}

func Test_Processor_ServeHTTP(t *testing.T) {
	const TestUrl = "http://example.org"
	aTest := tester.New(t)
	var err error
	var at hh.AverageTest
	var p *Processor
	var ps *ProcessorSettings
	var headersIn, headersOut http.Header
	var funcName string

	// Test #1. 'init' fails.
	{
		ps = &ProcessorSettings{}
		p, err = NewProcessor(ps)
		aTest.MustBeNoError(err)

		headersIn = http.Header{}
		headersOut = http.Header{}
		at = hh.AverageTest{
			Parameter: hh.AverageTestParameter{
				RequestMethod:  http.MethodGet,
				RequestUrl:     TestUrl,
				RequestHeaders: headersIn,
				RequestBody:    bytes.NewReader([]byte{}),
				RequestHandler: p.ServeHTTP,
			},
			ResultExpected: hh.AverageTestResult{
				ResponseStatusCode: http.StatusMethodNotAllowed,
				ResponseHeaders:    headersOut,
				ResponseBody:       []byte{},
			},
		}

		err = hh.PerformAverageHttpTest(&at)
		aTest.MustBeNoError(err)
		aTest.MustBeEqual(at.ResultReceived, at.ResultExpected)
	}

	// Test #2. 'run' fails.
	{
		ps = &ProcessorSettings{
			RequestIdFieldName: &_metaFieldName_RID,
		}
		p, err = NewProcessor(ps)
		aTest.MustBeNoError(err)
		err = p.AddFunc(RpcFunctionExampleTwo)
		aTest.MustBeNoError(err)

		headersIn = http.Header{}
		headersIn.Add(header.HttpHeaderContentType, mime.TypeApplicationJson)
		headersIn.Add(header.HttpHeaderAccept, mime.TypeAny)
		headersOut = http.Header{}
		headersOut.Add(header.HttpHeaderContentType, mime.TypeApplicationJson)
		funcName = "RpcFunctionExampleTwo"
		reqStr := fmt.Sprintf(`{"jsonrpc":"M1","id":"2","method":"%s","params":{}}`, funcName)
		at = hh.AverageTest{
			Parameter: hh.AverageTestParameter{
				RequestMethod:  http.MethodPost,
				RequestUrl:     TestUrl,
				RequestHeaders: headersIn,
				RequestBody:    strings.NewReader(reqStr),
				RequestHandler: p.ServeHTTP,
			},
			ResultExpected: hh.AverageTestResult{
				ResponseStatusCode: http.StatusOK,
				ResponseHeaders:    headersOut,
				ResponseBody:       []byte(`{"jsonrpc":"M1","id":"2","result":null,"error":{"code":-32,"message":"Internal RPC error","data":null},"ok":false}` + "\n"),
			},
		}

		err = hh.PerformAverageHttpTest(&at)
		aTest.MustBeNoError(err)
		aTest.MustBeEqual(at.ResultReceived, at.ResultExpected)
	}

	// Test #3. 'respond'.
	{
		ps = &ProcessorSettings{
			RequestIdFieldName: &_metaFieldName_RID,
		}
		p, err = NewProcessor(ps)
		aTest.MustBeNoError(err)
		err = p.AddFunc(RpcFunctionSum)
		aTest.MustBeNoError(err)

		headersIn = http.Header{}
		headersIn.Add(header.HttpHeaderContentType, mime.TypeApplicationJson)
		headersIn.Add(header.HttpHeaderAccept, mime.TypeAny)
		headersOut = http.Header{}
		headersOut.Add(header.HttpHeaderContentType, mime.TypeApplicationJson)
		funcName = "RpcFunctionSum"
		reqStr := fmt.Sprintf(`{"jsonrpc":"M1","id":"3","method":"%s","params":{"a":1,"b":2}}`, funcName)
		at = hh.AverageTest{
			Parameter: hh.AverageTestParameter{
				RequestMethod:  http.MethodPost,
				RequestUrl:     TestUrl,
				RequestHeaders: headersIn,
				RequestBody:    strings.NewReader(reqStr),
				RequestHandler: p.ServeHTTP,
			},
			ResultExpected: hh.AverageTestResult{
				ResponseStatusCode: http.StatusOK,
				ResponseHeaders:    headersOut,
				ResponseBody:       []byte(`{"jsonrpc":"M1","id":"3","result":{"c":3},"error":null,"ok":true}` + "\n"),
			},
		}

		err = hh.PerformAverageHttpTest(&at)
		aTest.MustBeNoError(err)
		aTest.MustBeEqual(at.ResultReceived, at.ResultExpected)
	}

}

func Test_Processor_GetRequestsCount(t *testing.T) {
	aTest := tester.New(t)
	ps := &ProcessorSettings{}
	p, err := NewProcessor(ps)
	aTest.MustBeNoError(err)
	a, b := p.GetRequestsCount()
	aTest.MustBeEqual(a, "0")
	aTest.MustBeEqual(b, "0")
}

func Test_Processor_incAllRequestsCounter(t *testing.T) {
	aTest := tester.New(t)
	ps := &ProcessorSettings{
		CountRequests: true,
	}
	p, err := NewProcessor(ps)
	aTest.MustBeNoError(err)
	p.incAllRequestsCounter()
	a, b := p.GetRequestsCount()
	aTest.MustBeEqual(a, "1")
	aTest.MustBeEqual(b, "0")
}

func Test_Processor_incSuccessfulRequestsCounter(t *testing.T) {
	aTest := tester.New(t)
	ps := &ProcessorSettings{
		CountRequests: true,
	}
	p, err := NewProcessor(ps)
	aTest.MustBeNoError(err)
	p.incSuccessfulRequestsCounter()
	a, b := p.GetRequestsCount()
	aTest.MustBeEqual(a, "0")
	aTest.MustBeEqual(b, "1")
}
