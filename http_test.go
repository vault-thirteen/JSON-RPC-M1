package jrm1

import (
	"bytes"
	"net/http"
	"testing"

	mime "github.com/vault-thirteen/auxie/MIME"
	"github.com/vault-thirteen/auxie/header"
	hh "github.com/vault-thirteen/auxie/http-helper"
	"github.com/vault-thirteen/auxie/tester"
)

func Test_checkHttpRequest(t *testing.T) {
	const TestUrl = "http://example.org"
	aTest := tester.New(t)
	var err error
	var proceed bool
	var headersIn, headersOut http.Header
	var at hh.AverageTest

	// 1.
	{
		headersIn = http.Header{}
		headersOut = http.Header{}
		at = hh.AverageTest{
			Parameter: hh.AverageTestParameter{
				RequestMethod:  http.MethodGet,
				RequestUrl:     TestUrl,
				RequestHeaders: headersIn,
				RequestBody:    bytes.NewReader([]byte{}),
				RequestHandler: func(rw http.ResponseWriter, req *http.Request) {
					proceed = checkHttpRequest(rw, req)
				},
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
		aTest.MustBeEqual(proceed, false)
	}

	// 2.
	{
		headersIn = http.Header{}
		headersIn.Add(header.HttpHeaderContentType, mime.TypeTextPlain)
		headersOut = http.Header{}
		at = hh.AverageTest{
			Parameter: hh.AverageTestParameter{
				RequestMethod:  http.MethodPost,
				RequestUrl:     TestUrl,
				RequestHeaders: headersIn,
				RequestBody:    bytes.NewReader([]byte{}),
				RequestHandler: func(rw http.ResponseWriter, req *http.Request) {
					proceed = checkHttpRequest(rw, req)
				},
			},
			ResultExpected: hh.AverageTestResult{
				ResponseStatusCode: http.StatusUnsupportedMediaType,
				ResponseHeaders:    headersOut,
				ResponseBody:       []byte{},
			},
		}

		err = hh.PerformAverageHttpTest(&at)
		aTest.MustBeNoError(err)
		aTest.MustBeEqual(at.ResultReceived, at.ResultExpected)
		aTest.MustBeEqual(proceed, false)
	}

	// 3.
	{
		headersIn = http.Header{}
		headersIn.Add(header.HttpHeaderContentType, mime.TypeApplicationJson)
		headersIn.Add(header.HttpHeaderAccept, mime.TypeTextPlain)
		headersOut = http.Header{}
		at = hh.AverageTest{
			Parameter: hh.AverageTestParameter{
				RequestMethod:  http.MethodPost,
				RequestUrl:     TestUrl,
				RequestHeaders: headersIn,
				RequestBody:    bytes.NewReader([]byte{}),
				RequestHandler: func(rw http.ResponseWriter, req *http.Request) {
					proceed = checkHttpRequest(rw, req)
				},
			},
			ResultExpected: hh.AverageTestResult{
				ResponseStatusCode: http.StatusNotAcceptable,
				ResponseHeaders:    headersOut,
				ResponseBody:       []byte{},
			},
		}

		err = hh.PerformAverageHttpTest(&at)
		aTest.MustBeNoError(err)
		aTest.MustBeEqual(at.ResultReceived, at.ResultExpected)
		aTest.MustBeEqual(proceed, false)
	}

	// 4.
	{
		headersIn = http.Header{}
		headersIn.Add(header.HttpHeaderContentType, mime.TypeApplicationJson)
		headersIn.Add(header.HttpHeaderAccept, mime.TypeAny)
		headersOut = http.Header{}
		at = hh.AverageTest{
			Parameter: hh.AverageTestParameter{
				RequestMethod:  http.MethodPost,
				RequestUrl:     TestUrl,
				RequestHeaders: headersIn,
				RequestBody:    bytes.NewReader([]byte{}),
				RequestHandler: func(rw http.ResponseWriter, req *http.Request) {
					proceed = checkHttpRequest(rw, req)
				},
			},
			ResultExpected: hh.AverageTestResult{
				ResponseStatusCode: http.StatusOK,
				ResponseHeaders:    headersOut,
				ResponseBody:       []byte{},
			},
		}

		err = hh.PerformAverageHttpTest(&at)
		aTest.MustBeNoError(err)
		aTest.MustBeEqual(at.ResultReceived, at.ResultExpected)
		aTest.MustBeEqual(proceed, true)
	}
}
