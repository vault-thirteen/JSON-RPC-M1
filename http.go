package jrm1

import (
	"net/http"

	mime "github.com/vault-thirteen/auxie/MIME"
	h "github.com/vault-thirteen/auxie/header"
	hh "github.com/vault-thirteen/auxie/http-helper"
)

// checkHttpRequest checks the HTTP request and responds on error.
// If the request is correct and ready to be processed, 'True' is returned.
// When 'False' is returned, the caller must stop serving the request.
func checkHttpRequest(rw http.ResponseWriter, req *http.Request) (proceed bool) {
	// 1. Check HTTP method.
	if req.Method != http.MethodPost {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return false
	}

	// 2. Check HTTP content type.
	ctype, err := hh.GetSingleHttpHeader(req, h.HttpHeaderContentType)
	if (err != nil) || (ctype != mime.TypeApplicationJson) {
		rw.WriteHeader(http.StatusUnsupportedMediaType)
		return false
	}

	// 3. Check accepted HTTP content type.
	var ok bool
	ok, err = hh.CheckBrowserSupportForJson(req)
	if (err != nil) || !ok {
		rw.WriteHeader(http.StatusNotAcceptable)
		return false
	}

	return true
}
