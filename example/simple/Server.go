package main

import (
	"net/http"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
)

const (
	DurationFieldName  = "dur"
	RequestIdFieldName = "rid"
)

type Server struct {
	p *jrm1.Processor
}

func NewServer() (s *Server, err error) {
	s = &Server{}

	durationFieldName := DurationFieldName
	requestIdFieldName := RequestIdFieldName

	ps := &jrm1.ProcessorSettings{
		CatchExceptions:    true,
		LogExceptions:      true,
		CountRequests:      true,
		DurationFieldName:  &durationFieldName,
		RequestIdFieldName: &requestIdFieldName,
	}

	s.p, err = jrm1.NewProcessor(ps)
	if err != nil {
		return nil, err
	}

	s.p.AddFuncFast(Sum)
	s.p.AddFuncFast(Crash)

	return s, nil
}

func (s *Server) rootHandler(rw http.ResponseWriter, req *http.Request) {
	s.p.ServeHTTP(rw, req)
}
