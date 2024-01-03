package jrm1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"sync"

	mime "github.com/vault-thirteen/auxie/MIME"
	ae "github.com/vault-thirteen/auxie/errors"
	"github.com/vault-thirteen/auxie/header"
)

// Client is an RPC client.
type Client struct {
	settings *ClientSettings
	guard    *sync.RWMutex

	// Request counters.
	requestsCount    *big.Int
	requestsCountOne *big.Int
}

// NewClient creates an RPC client.
func NewClient(settings *ClientSettings) (c *Client, err error) {
	err = settings.Check()
	if err != nil {
		return nil, err
	}

	c = &Client{
		settings:         settings,
		guard:            new(sync.RWMutex),
		requestsCount:    big.NewInt(0),
		requestsCountOne: big.NewInt(1),
	}

	return c, nil
}

// Call performs a function call and puts result into the 'result' argument.
// The 'result' argument must be a pointer to an initialised (empty) object.
func (c *Client) Call(ctx context.Context, method string, params any, result any) (re *RpcError, err error) {
	c.guard.Lock()
	defer c.guard.Unlock()

	// Prepare protocol name and request ID.
	pn := ProtocolNameM1
	c.incRequestsCount()
	var rid = c.GetRequestsCount()

	// Encode parameters.
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "\t")
	enc.SetEscapeHTML(c.settings.useHtmlEscaping)
	err = enc.Encode(params)
	if err != nil {
		return nil, err
	}
	var paramsJRM = json.RawMessage(buf.Bytes())

	// Make a function call.
	rpcReq := &RpcRequest{
		ProtocolName: &pn,
		Id:           &rid,
		Method:       &method,
		Parameters:   &paramsJRM,
	}

	var rawRpcResp *RpcResponseRaw
	rawRpcResp, err = c.call(ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// Get response error.
	re = rawRpcResp.Error

	// Get response result.
	if rawRpcResp.Result != nil {
		decoder := json.NewDecoder(bytes.NewReader(*rawRpcResp.Result))
		decoder.DisallowUnknownFields()
		decoder.UseNumber()
		err = decoder.Decode(result)
		if err != nil {
			return nil, err
		}
	}

	return re, nil
}

// CallRaw takes a raw request, performs the request, returns a raw response.
func (c *Client) CallRaw(ctx context.Context, rpcReq *RpcRequest) (rpcResp *RpcResponseRaw, err error) {
	c.guard.Lock()
	defer c.guard.Unlock()

	return c.call(ctx, rpcReq)
}

// call makes a request to the RPC server, receives a response and returns it.
// Result object is set as an empty interface ('any' type).
func (c *Client) call(ctx context.Context, rpcReq *RpcRequest) (rpcResp *RpcResponseRaw, err error) {
	if !rpcReq.HasAllRootFields() {
		return nil, errors.New(ErrRpcRequestIsMalformed)
	}

	err = rpcReq.CheckProtocolVersion()
	if err != nil {
		return nil, err
	}

	var httpReq *http.Request
	httpReq, err = c.newHttpRequest(ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	var httpClient *http.Client
	if c.settings.httpClient != nil {
		httpClient = c.settings.httpClient
	} else {
		httpClient = new(http.Client)
	}

	var httpResp *http.Response
	httpResp, err = httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer func() {
		derr := httpResp.Body.Close()
		if derr != nil {
			err = ae.Combine(err, derr)
		}
	}()

	decoder := json.NewDecoder(httpResp.Body)
	decoder.DisallowUnknownFields()
	decoder.UseNumber()
	err = decoder.Decode(&rpcResp)
	if err != nil {
		return nil, err
	}

	return rpcResp, nil
}

// newHttpRequest takes an RPC request object and creates an HTTP request for
// the RPC server.
func (c *Client) newHttpRequest(ctx context.Context, rpcReq *RpcRequest) (hr *http.Request, err error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "\t")
	enc.SetEscapeHTML(c.settings.useHtmlEscaping)
	err = enc.Encode(rpcReq)
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf(
		"%s://%s:%d%s",
		c.settings.schema,
		c.settings.host,
		c.settings.port,
		c.settings.path,
	)

	hr, err = http.NewRequestWithContext(ctx, http.MethodPost, dsn, &buf)
	if err != nil {
		return nil, err
	}

	hr.Header.Set(header.HttpHeaderContentType, mime.TypeApplicationJson)
	hr.Header.Set(header.HttpHeaderAccept, mime.TypeApplicationJson)

	// Set the custom HTTP headers for those who need them.
	for hdrName, hdrValue := range c.settings.httpHeaders {
		hr.Header.Set(hdrName, hdrValue)
	}

	return hr, nil
}

// GetRequestsCount returns the counter of performed calls (requests) to the RPC
// server.
func (c *Client) GetRequestsCount() (requestsCount string) {
	return c.requestsCount.String()
}

// incRequestsCount increases the counter of performed calls (requests) to the
// RPC server by one.
func (c *Client) incRequestsCount() {
	c.requestsCount.Add(c.requestsCount, c.requestsCountOne)
}
