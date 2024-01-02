package jrm1

import (
	"errors"
	"net/http"
)

const (
	ErrClientSettingsError = "error is client settings"
)

// ClientSettings are settings of an RPC client.
type ClientSettings struct {
	// Target server URL schema.
	schema string

	// Target server host name.
	host string

	// Target server port number.
	port uint16

	// Relative URL path to which the client makes requests.
	// A typical URL path is "/", however it may be different if a custom RPC
	// server is used.
	path string

	// Custom HTTP client.
	httpClient *http.Client

	// Custom HTTP headers.
	httpHeaders map[string]string

	// If enabled, some of HTML entities will be escaped during JSON encoding.
	useHtmlEscaping bool
}

// NewClientSettings is a constructor of an RPC client settings.
func NewClientSettings(schema string, host string, port uint16, path string, customHttpClient *http.Client, customHttpHeaders map[string]string, useHtmlEscaping bool) (cs *ClientSettings, err error) {
	cs = &ClientSettings{
		schema:          schema,
		host:            host,
		port:            port,
		path:            path,
		httpClient:      customHttpClient,
		httpHeaders:     customHttpHeaders,
		useHtmlEscaping: useHtmlEscaping,
	}

	err = cs.Check()
	if err != nil {
		return nil, err
	}

	return cs, nil
}

// Check validates the settings of an RPC client.
func (cs *ClientSettings) Check() (err error) {
	if (len(cs.schema) == 0) ||
		(len(cs.host) == 0) ||
		(cs.port == 0) ||
		(len(cs.path) == 0) {
		return errors.New(ErrClientSettingsError)
	}

	return nil
}
