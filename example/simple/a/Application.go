package a

import (
	"github.com/vault-thirteen/JSON-RPC-M1/example/simple/s"
)

type Application struct {
	s *s.Server
}

func NewApplication() (a *Application, err error) {
	a = new(Application)

	a.s, err = s.NewServer()
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *Application) Start() {
	a.s.Start()
}

func (a *Application) Stop() (err error) {
	err = a.s.Stop()
	if err != nil {
		return err
	}

	return nil
}
