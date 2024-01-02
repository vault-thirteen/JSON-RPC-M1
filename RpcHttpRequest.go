package jrm1

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	mime "github.com/vault-thirteen/auxie/MIME"
	"github.com/vault-thirteen/auxie/header"
)

// RpcHttpRequest is an RPC request originated from an HTTP request.
type RpcHttpRequest struct {
	// Processor.
	p *Processor

	// Processor settings.
	settings *ProcessorSettings

	// HTTP request from which the RPC request was received.
	req *http.Request

	// HTTP response writer used for communication with client.
	rw http.ResponseWriter

	// Time of start of the request processing.
	tStart time.Time

	// Time spent on processing a successful function call.
	// Is set in milliseconds.
	tDuration int64

	// RPC request.
	rr *RpcRequest

	// RPC response.
	resp *RpcResponse
}

// NewRpcHttpRequest is a simple constructor of an RPC request originated from
// an HTTP request.
func NewRpcHttpRequest(p *Processor, settings *ProcessorSettings, req *http.Request, rw http.ResponseWriter) (rhr *RpcHttpRequest) {
	rhr = &RpcHttpRequest{
		p:        p,
		settings: settings,
		req:      req,
		rw:       rw,
	}

	return rhr
}

// init reads a request from HTTP body, checks it and searches for the
// requested function. If error occurs, it responds to the client via HTTP.
// If request is correct and ready to be processed further, 'True' is returned.
// When 'False' is returned, the caller must stop serving the request.
func (r *RpcHttpRequest) init() (proceed bool) {
	r.startTimer()

	r.p.incAllRequestsCounter()

	if !checkHttpRequest(r.rw, r.req) {
		return false
	}

	r.resp = NewRpcResponse()

	var err error
	r.rr, err = NewRpcRequest(r.req.Body)
	if err != nil {
		r.resp.Error = NewRpcErrorFast(RpcErrorCode_RequestIsNotReadable)
		r.respond()
		return false
	}

	r.resp.Id = r.rr.Id

	var ok bool
	ok = r.rr.HasAllRootFields()
	if !ok {
		r.resp.Error = NewRpcErrorFast(RpcErrorCode_InvalidRequest)
		r.respond()
		return false
	}

	err = r.rr.CheckProtocolVersion()
	if err != nil {
		r.resp.Error = NewRpcErrorFast(RpcErrorCode_UnsupportedProtocol)
		r.respond()
		return false
	}

	err = r.p.FindFunc(*r.rr.Method)
	if err != nil {
		r.resp.Error = NewRpcErrorFast(RpcErrorCode_UnknownMethod)
		r.respond()
		return false
	}

	return true
}

// startTimer starts the timer.
func (r *RpcHttpRequest) startTimer() {
	if r.settings.isDurationEnabled() {
		r.tStart = time.Now()
	}
}

// run calls the requested function and collects the result. If error occurs,
// it responds to the client via HTTP. If everything is correct and ready to be
// processed further, 'True' is returned. When 'False' is returned, the caller
// must stop serving the request.
func (r *RpcHttpRequest) run() (proceed bool) {
	var err error
	if r.settings.isRequestIdShown() {
		err = r.resp.Meta.AddField(*r.settings.RequestIdFieldName, *r.rr.Id)
		if err != nil {
			r.resp.Error = NewRpcErrorFast(RpcErrorCode_InternalRpcError)
			r.respond()
			return false
		}
	}

	r.resp.Result, r.resp.Error = r.p.RunFunc(*r.rr.Method, r.rr.Parameters, r.resp.Meta)

	if r.settings.isRequestIdShown() {
		err = r.resp.Meta.RemoveField(*r.settings.RequestIdFieldName)
		if err != nil {
			r.resp.Error = NewRpcErrorFast(RpcErrorCode_InternalRpcError)
			r.respond()
			return false
		}
	}

	if !r.stopTimer() {
		return false
	}

	return true
}

// stopTimer stops the timer and saves the duration as a meta-data field. If
// error occurs, it responds to the client via HTTP. If everything is correct
// and ready to be processed further, 'True' is returned. When 'False' is
// returned, the caller must stop serving the request.
func (r *RpcHttpRequest) stopTimer() (proceed bool) {
	if r.settings.isDurationEnabled() {
		r.tDuration = time.Now().Sub(r.tStart).Milliseconds()

		err := r.resp.Meta.AddField(*r.p.settings.DurationFieldName, r.tDuration)
		if err != nil {
			r.resp.Error = NewRpcErrorFast(RpcErrorCode_InternalRpcError)
			r.respond()
			return false
		}
	}

	return true
}

// respond analyses the result and responds to the client via HTTP. The caller
// must stop serving the request after this function returns.
func (r *RpcHttpRequest) respond() {
	// Empty meta-data set must not be shown.
	if len(*r.resp.Meta) == 0 {
		r.resp.Meta = nil
	}

	if !r.resp.hasError() {
		r.resp.OK = true
		r.p.incSuccessfulRequestsCounter()
	}

	r.rw.Header().Set(header.HttpHeaderContentType, mime.TypeApplicationJson)

	err := json.NewEncoder(r.rw).Encode(r.resp)
	if err != nil {
		log.Println(err)
	}
}
