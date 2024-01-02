package jrm1

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	mime "github.com/vault-thirteen/auxie/MIME"
	"github.com/vault-thirteen/auxie/header"
	"github.com/vault-thirteen/auxie/tester"
)

func Test_NewRpcHttpRequest(t *testing.T) {
	aTest := tester.New(t)
	var rhr, rhrExpected *RpcHttpRequest

	var (
		p        *Processor
		settings *ProcessorSettings
		req      *http.Request
		rw       http.ResponseWriter
	)

	// Test.
	p = new(Processor)
	settings = new(ProcessorSettings)
	req = new(http.Request)
	rw = new(httptest.ResponseRecorder)

	rhrExpected = &RpcHttpRequest{
		p:        p,
		settings: settings,
		req:      req,
		rw:       rw,
	}

	rhr = NewRpcHttpRequest(p, settings, req, rw)
	aTest.MustBeEqual(rhr, rhrExpected)
}

func Test_RpcHttpRequest_init(t *testing.T) {
	aTest := tester.New(t)
	var r *RpcHttpRequest
	var p *Processor
	var ps *ProcessorSettings
	var err error
	var recorder http.ResponseWriter
	var req *http.Request
	var proceed bool

	// Test #1. 'checkHttpRequest' fails.
	{
		ps = &ProcessorSettings{}
		p, err = NewProcessor(ps)
		aTest.MustBeNoError(err)
		recorder = httptest.NewRecorder()
		req = &http.Request{}
		req.Method = http.MethodGet
		r = NewRpcHttpRequest(p, ps, req, recorder)
		proceed = r.init()
		//
		aTest.MustBeEqual(proceed, false)
	}

	// Test #2. 'NewRpcRequest' fails.
	{
		ps = &ProcessorSettings{}
		p, err = NewProcessor(ps)
		aTest.MustBeNoError(err)
		recorder = httptest.NewRecorder()
		req = &http.Request{}
		req.Method = http.MethodPost
		req.Header = http.Header{}
		req.Header.Add(header.HttpHeaderContentType, mime.TypeApplicationJson)
		req.Header.Add(header.HttpHeaderAccept, mime.TypeAny)
		req.Body = io.NopCloser(strings.NewReader("Not a JSON, haha !"))
		r = NewRpcHttpRequest(p, ps, req, recorder)
		proceed = r.init()
		//
		aTest.MustBeEqual(proceed, false)
	}

	// Test #3. 'HasAllRootFields' fails.
	{
		ps = &ProcessorSettings{}
		p, err = NewProcessor(ps)
		aTest.MustBeNoError(err)
		recorder = httptest.NewRecorder()
		req = &http.Request{}
		req.Method = http.MethodPost
		req.Header = http.Header{}
		req.Header.Add(header.HttpHeaderContentType, mime.TypeApplicationJson)
		req.Header.Add(header.HttpHeaderAccept, mime.TypeAny)
		req.Body = io.NopCloser(strings.NewReader(`{"jsonrpc":"M1","method":"abc"}`)) // Not all root fields are set.
		r = NewRpcHttpRequest(p, ps, req, recorder)
		proceed = r.init()
		//
		aTest.MustBeEqual(proceed, false)
	}

	// Test #4. 'CheckProtocolVersion' fails.
	{
		ps = &ProcessorSettings{}
		p, err = NewProcessor(ps)
		aTest.MustBeNoError(err)
		recorder = httptest.NewRecorder()
		req = &http.Request{}
		req.Method = http.MethodPost
		req.Header = http.Header{}
		req.Header.Add(header.HttpHeaderContentType, mime.TypeApplicationJson)
		req.Header.Add(header.HttpHeaderAccept, mime.TypeAny)
		req.Body = io.NopCloser(strings.NewReader(`{"jsonrpc":"???","id":"1","method":"m","params":{}}`))
		r = NewRpcHttpRequest(p, ps, req, recorder)
		proceed = r.init()
		//
		aTest.MustBeEqual(proceed, false)
	}

	// Test #5. 'FindFunc' fails.
	{
		ps = &ProcessorSettings{}
		p, err = NewProcessor(ps)
		aTest.MustBeNoError(err)
		recorder = httptest.NewRecorder()
		req = &http.Request{}
		req.Method = http.MethodPost
		req.Header = http.Header{}
		req.Header.Add(header.HttpHeaderContentType, mime.TypeApplicationJson)
		req.Header.Add(header.HttpHeaderAccept, mime.TypeAny)
		req.Body = io.NopCloser(strings.NewReader(`{"jsonrpc":"M1","id":"1","method":"x","params":{}}`))
		r = NewRpcHttpRequest(p, ps, req, recorder)
		proceed = r.init()
		//
		aTest.MustBeEqual(proceed, false)
	}

	// Test #6. OK.
	{
		ps = &ProcessorSettings{}
		p, err = NewProcessor(ps)
		aTest.MustBeNoError(err)
		err = p.AddFunc(RpcFunctionExampleOne)
		aTest.MustBeNoError(err)
		recorder = httptest.NewRecorder()
		req = &http.Request{}
		req.Method = http.MethodPost
		req.Header = http.Header{}
		req.Header.Add(header.HttpHeaderContentType, mime.TypeApplicationJson)
		req.Header.Add(header.HttpHeaderAccept, mime.TypeAny)
		req.Body = io.NopCloser(strings.NewReader(`{"jsonrpc":"M1","id":"1","method":"RpcFunctionExampleOne","params":{}}`))
		r = NewRpcHttpRequest(p, ps, req, recorder)
		proceed = r.init()
		//
		aTest.MustBeEqual(proceed, true)
	}
}

func Test_RpcHttpRequest_startTimer(t *testing.T) {
	aTest := tester.New(t)
	var r *RpcHttpRequest
	var p *Processor
	var ps *ProcessorSettings
	var err error

	// Test #1. Duration is enabled.
	durFN := "dur"
	ps = &ProcessorSettings{DurationFieldName: &durFN}
	p, err = NewProcessor(ps)
	aTest.MustBeNoError(err)
	r = NewRpcHttpRequest(p, ps, nil, nil)
	aTest.MustBeEqual(r.tStart, time.Time{})
	r.startTimer()
	aTest.MustBeDifferent(r.tStart, time.Time{})

	// Test #2. Duration is disabled.
	ps = &ProcessorSettings{DurationFieldName: nil}
	p, err = NewProcessor(ps)
	aTest.MustBeNoError(err)
	r = NewRpcHttpRequest(p, ps, nil, nil)
	aTest.MustBeEqual(r.tStart, time.Time{})
	r.startTimer()
	aTest.MustBeEqual(r.tStart, time.Time{})
}

func Test_RpcHttpRequest_run(t *testing.T) {
	aTest := tester.New(t)
	var r *RpcHttpRequest
	var p *Processor
	var ps *ProcessorSettings
	var err error
	var recorder http.ResponseWriter
	var req *http.Request
	var proceed bool

	testPreparer := func(showReqId bool, calcDuration bool, funcName string) {
		ps = &ProcessorSettings{}
		if showReqId {
			ps.RequestIdFieldName = &_metaFieldName_RID
		}
		if calcDuration {
			ps.DurationFieldName = &_metaFieldName_Duration
		}

		p, err = NewProcessor(ps)
		aTest.MustBeNoError(err)
		err = p.AddFunc(RpcFunctionExampleOne)
		aTest.MustBeNoError(err)
		err = p.AddFunc(RpcFunctionExampleTwo)
		aTest.MustBeNoError(err)
		err = p.AddFunc(RpcFunctionExampleThree)
		aTest.MustBeNoError(err)

		recorder = httptest.NewRecorder()
		req = &http.Request{}
		req.Method = http.MethodPost
		req.Header = http.Header{}
		req.Header.Add(header.HttpHeaderContentType, mime.TypeApplicationJson)
		req.Header.Add(header.HttpHeaderAccept, mime.TypeAny)
		reqStr := fmt.Sprintf(`{"jsonrpc":"M1","id":"1","method":"%s","params":{}}`, funcName)
		req.Body = io.NopCloser(strings.NewReader(reqStr))
		r = NewRpcHttpRequest(p, ps, req, recorder)
		proceed = r.init()
		aTest.MustBeEqual(proceed, true)
	}

	// Test #1. Request ID is shown, meta-data fails.
	testPreparer(true, false, "RpcFunctionExampleOne")
	err = r.resp.Meta.AddField(_metaFieldName_RID, "junk") // Duplicate RID field.
	aTest.MustBeNoError(err)
	proceed = r.run()
	aTest.MustBeEqual(proceed, false)

	// Test #2. Request ID is shown, meta-data fails after function call.
	testPreparer(true, false, "RpcFunctionExampleTwo")
	proceed = r.run()
	aTest.MustBeEqual(proceed, false)

	// Test #3. Request ID is shown, 'stopTimer' fails.
	testPreparer(true, true, "RpcFunctionExampleThree")
	proceed = r.run()
	aTest.MustBeEqual(proceed, false)

	// Test #3. Request ID is shown, 'stopTimer' fails.
	testPreparer(true, false, "RpcFunctionExampleOne")
	proceed = r.run()
	aTest.MustBeEqual(proceed, true)
}

func Test_RpcHttpRequest_stopTimer(t *testing.T) {
	aTest := tester.New(t)
	var r *RpcHttpRequest
	var p *Processor
	var ps *ProcessorSettings
	var err error
	var proceed bool
	durFN := "dur"

	// Test #1. Duration is disabled.
	ps = &ProcessorSettings{DurationFieldName: nil}
	p, err = NewProcessor(ps)
	aTest.MustBeNoError(err)
	r = NewRpcHttpRequest(p, ps, nil, nil)
	//
	aTest.MustBeEqual(r.tStart, time.Time{})
	proceed = r.stopTimer()
	aTest.MustBeEqual(proceed, true)
	aTest.MustBeEqual(r.tDuration, int64(0))

	// Test #2. Duration is enabled, no error.
	ps = &ProcessorSettings{DurationFieldName: &durFN}
	p, err = NewProcessor(ps)
	aTest.MustBeNoError(err)
	r = NewRpcHttpRequest(p, ps, nil, nil)
	r.resp = NewRpcResponse()
	//
	r.startTimer()
	time.Sleep(time.Second)
	proceed = r.stopTimer()
	aTest.MustBeEqual(proceed, true)
	aTest.MustBeDifferent(r.tDuration, int64(0))

	// Test #3. Duration is enabled, duplicate meta-data field.
	ps = &ProcessorSettings{DurationFieldName: &durFN}
	p, err = NewProcessor(ps)
	aTest.MustBeNoError(err)
	r = NewRpcHttpRequest(p, ps, nil, httptest.NewRecorder())
	r.resp = NewRpcResponse()
	err = r.resp.Meta.AddField(durFN, "duplicate key value")
	aTest.MustBeNoError(err)
	//
	r.startTimer()
	time.Sleep(time.Second)
	proceed = r.stopTimer()
	aTest.MustBeEqual(proceed, false)
	aTest.MustBeEqual(r.tDuration > int64(1000), true) // ms.
}

func Test_RpcHttpRequest_respond(t *testing.T) {
	aTest := tester.New(t)
	var r *RpcHttpRequest
	var p *Processor
	var ps *ProcessorSettings
	var err error
	var respBody []byte

	// Test #1. Error is not set.
	ps = &ProcessorSettings{}
	p, err = NewProcessor(ps)
	aTest.MustBeNoError(err)
	recorder := httptest.NewRecorder()
	r = NewRpcHttpRequest(p, ps, nil, recorder)
	r.resp = NewRpcResponse()
	r.resp.Error = nil
	r.respond()
	//
	aTest.MustBeEqual(r.resp.OK, true)
	aTest.MustBeEqual(recorder.Header().Get(header.HttpHeaderContentType), mime.TypeApplicationJson)
	respBody, err = io.ReadAll(recorder.Body)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(
		strings.TrimSpace(string(respBody)),
		`{"jsonrpc":"M1","id":null,"result":null,"error":null,"ok":true}`,
	)
}
