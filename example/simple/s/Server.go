package s

import (
	"context"
	"log"
	"net/http"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
)

const (
	DurationFieldName  = "dur"
	RequestIdFieldName = "rid"
)

type Server struct {
	p  *jrm1.Processor
	hs *http.Server
}

func NewServer() (s *Server, err error) {
	s = new(Server)

	err = s.initProcessor()
	if err != nil {
		return nil, err
	}

	err = s.initHttpServer()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Server) initProcessor() (err error) {
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
		return err
	}

	s.p.AddFuncFast(Sum)
	s.p.AddFuncFast(Crash)

	return nil
}

func (s *Server) initHttpServer() (err error) {
	s.hs = &http.Server{
		Handler: http.HandlerFunc(s.RootHandler),
	}

	return nil
}

func (s *Server) RootHandler(rw http.ResponseWriter, req *http.Request) {
	s.p.ServeHTTP(rw, req)
}

func (s *Server) Start() {
	go s.run()
}

func (s *Server) run() {
	err := s.hs.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}

func (s *Server) Stop() (err error) {
	err = s.hs.Shutdown(context.Background())
	if err != nil {
		return err
	}

	return nil
}
