package jrm1

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"runtime/debug"
	"sync"
)

const (
	ErrDuplicateFunction  = "duplicate function"
	ErrFunctionIsNotFound = "function is not found"
)

// Processor is an RPC processor (server).
type Processor struct {
	settings *ProcessorSettings
	guard    *sync.RWMutex

	// List of RPC functions.
	funcs map[string]RpcFunction

	// Request counters.
	requestsCountAll        *big.Int
	requestsCountSuccessful *big.Int
	requestsCountOne        *big.Int
}

// NewProcessor is a constructor of an empty RPC processor (server).
func NewProcessor(settings *ProcessorSettings) (p *Processor, err error) {
	err = settings.Check()
	if err != nil {
		return nil, err
	}

	p = &Processor{
		settings:                settings,
		guard:                   new(sync.RWMutex),
		funcs:                   make(map[string]RpcFunction),
		requestsCountAll:        big.NewInt(0),
		requestsCountSuccessful: big.NewInt(0),
		requestsCountOne:        big.NewInt(1),
	}

	if errorMessages == nil {
		initErrorMessages()
	}

	return p, nil
}

// AddFunc tries to add a function to the RPC processor (server).
func (p *Processor) AddFunc(f RpcFunction) (err error) {
	p.guard.Lock()
	defer p.guard.Unlock()

	funcName := f.GetName()

	err = CheckFunctionName(funcName)
	if err != nil {
		return err
	}

	_, alreadyExists := p.funcs[funcName]
	if alreadyExists {
		return errors.New(ErrDuplicateFunction)
	}

	p.funcs[funcName] = f

	return nil
}

// AddFuncFast tries to add a function to the RPC processor (server).
// It panics on error.
func (p *Processor) AddFuncFast(f RpcFunction) {
	err := p.AddFunc(f)
	if err != nil {
		panic(err)
	}
}

// RemoveFunc tries to remove a function from the RPC processor (server).
func (p *Processor) RemoveFunc(funcName string) (err error) {
	p.guard.Lock()
	defer p.guard.Unlock()

	_, exists := p.funcs[funcName]
	if !exists {
		return errors.New(ErrFunctionIsNotFound)
	}

	delete(p.funcs, funcName)

	return nil
}

// FindFunc checks presence of the function in the RPC processor (server).
func (p *Processor) FindFunc(funcName string) (err error) {
	p.guard.RLock()
	defer p.guard.RUnlock()

	_, exists := p.funcs[funcName]
	if !exists {
		return errors.New(ErrFunctionIsNotFound)
	}

	return nil
}

// RunFunc executes a function of the RPC processor (server) specified by its
// name. If enabled in settings, it also catches any exception (panic) which
// may happen during the function execution.
func (p *Processor) RunFunc(funcName string, params *json.RawMessage, metaData *ResponseMetaData) (result any, re *RpcError) {
	p.guard.RLock()
	defer p.guard.RUnlock()

	if p.settings.CatchExceptions {
		defer func() {
			x := recover()
			if x != nil {
				if p.settings.LogExceptions {
					log.Println(fmt.Sprintf("%v, %s", x, string(debug.Stack())))
				}

				re = NewRpcErrorFast(RpcErrorCode_InternalRpcError)
			}
		}()
	}

	f, ok := p.funcs[funcName]
	if !ok {
		return nil, NewRpcErrorFast(RpcErrorCode_UnknownMethod)
	}

	return f(params, metaData)
}

// ServeHTTP handles an HTTP request and responds to it.
// 'ServeHTTP' is a required method of the 'http.Handler' interface.
func (p *Processor) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rhr := NewRpcHttpRequest(p, p.settings, req, rw)

	if !rhr.init() {
		return
	}

	if !rhr.run() {
		return
	}

	rhr.respond()
	return
}

// GetRequestsCount returns the number of all (received) and successful
// function calls.
func (p *Processor) GetRequestsCount() (all, successful string) {
	return p.requestsCountAll.String(), p.requestsCountSuccessful.String()
}

// incAllRequestsCounter increments the counter of all (received) function
// calls if the request counter is enabled.
func (p *Processor) incAllRequestsCounter() {
	if p.settings.CountRequests {
		p.requestsCountAll.Add(p.requestsCountAll, p.requestsCountOne)
	}
}

// incSuccessfulRequestsCounter increments the counter of successful function
// calls if the request counter is enabled.
func (p *Processor) incSuccessfulRequestsCounter() {
	if p.settings.CountRequests {
		p.requestsCountSuccessful.Add(p.requestsCountSuccessful, p.requestsCountOne)
	}
}
