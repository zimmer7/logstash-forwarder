// Licensed to Elasticsearch under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package system

import (
	"fmt"
	"github.com/elasticsearch/kriterium/component/process"
	"github.com/elasticsearch/kriterium/errors"
	"github.com/elasticsearch/kriterium/panics"
	"os"
)

func StartRegistry(basepath string) (Registrar, error) {
	registrar, e := newRegistrar(basepath)
	if e != nil {
		return nil, e
	}
	go registrar.run()

	return registrar, nil
}

// ----------------------------------------------------------------------------
// registrar
// ----------------------------------------------------------------------------

// registrar implements system doc registrar functionality and supports the
// lsf/system.Registrar interface.
type registrar struct {
	controller process.Controller // active-command controller
	process    process.Process    // active-command process
	sighandler process.Process
	reg        *registry
	ui         chan req
}

func newRegistrar(basepath string) (*registrar, error) {
	registry, e := openRegistry(basepath)
	if e != nil {
		return nil, e
	}

	cnc := process.NewCmdCtl()
	sigHandler := cnc.CommandOnSignal(process.Stop, os.Interrupt, os.Kill)

	regisrar := &registrar{
		cnc.Controller(),
		cnc.Process(),
		sigHandler,
		registry,
		make(chan req, 12),
	}
	return regisrar, nil
}

// Launches a go-routine to process user requests affecting
// shared system resources managed by system, per semantics of system.Registrar.
func (r *registrar) run() {
	defer panics.AsyncRecover(r.controller.Respond())

	for {
		select {
		case request := <-r.ui:
			request.result <- request.execute()
		case cmd := <-r.controller.Command():
			var cmdresp process.CommandCode
			switch cmd {
			case process.Stop: // expected
				cmdresp = process.Stop
			default: // unexpected cmd.
				cmdresp = process.Abort
			}
			r.sighandler.Signal() <- process.Stop
			r.controller.Respond() <- cmdresp
			return
		}
	}
}

// ----------------------------------------------------------------------------
// interface: Registrar
// ----------------------------------------------------------------------------

func (r *registrar) Stop() (err error) {
	r.process.Signal() <- process.Stop
	if resp := <-r.process.Response(); resp != process.Stop {
		return errors.IllegalState("registrar.Stop:", "unexpected response:", resp)
	}
	return
}

func (r *registrar) String() string {
	s := fmt.Sprintf("registrar: path %s", r.reg.path)
	return s
}

func (r *registrar) DeleteDocument(id string) (bool, error) {
	fn := func() interface{} {
		ok, e := r.reg.deleteDocument(id)
		if e != nil {
			return e
		}
		return ok
	}
	return r.dispatch1(fn)
}

func (r *registrar) UpdateDocument(doc Document) (bool, error) {
	fn := func() interface{} {
		ok, e := r.reg.updateDocument(doc.(*document))
		if e != nil {
			return e
		}
		return ok
	}
	return r.dispatch1(fn)
}
func (r *registrar) ReadDocument(id string) (Document, error) {
	fn := func() interface{} {
		doc, e := r.reg.readDocument(id)
		if e != nil {
			return e
		}
		return doc
	}
	return r.dispatch0(fn)
}

func (r *registrar) CreateDocument(id string, data map[string][]byte) (Document, error) {
	fn := func() interface{} {
		doc, e := r.reg.createDocument(id, data)
		if e != nil {
			return e
		}
		return doc
	}
	return r.dispatch0(fn)
}

func (r *registrar) dispatch0(fn func() interface{}) (Document, error) {
	resch := makeResChan()
	r.ui <- req{resch, fn}
	result := <-resch
	return mapDocResult(result)
}

func (r *registrar) dispatch1(fn func() interface{}) (bool, error) {
	resch := makeResChan()
	r.ui <- req{resch, fn}
	result := <-resch
	return mapBoolResult(result)
}

// ----------------------------------------------------------------------------
// concurrent async request dispatch mechanism
// ----------------------------------------------------------------------------

type stat struct {
	err error
	dat []byte
}

// query type just wraps the delegated registrar func invokes
type query func() interface{}

// an async request is a tuple wrapping result callback channel
// and the actual (query) func invoke
type req struct {
	result  chan<- interface{}
	execute query
}

func makeResChan() chan interface{} { return make(chan interface{}, 1) }

func mapDocResult(result interface{}) (Document, error) {
	switch t := result.(type) {
	case Document:
		return t, nil
	case error:
		return nil, t
	default:
		panic("BUG - unexpected type value")
	}
}

func mapBoolResult(result interface{}) (bool, error) {
	switch t := result.(type) {
	case bool:
		return t, nil
	case error:
		return false, t
	default:
		panic("BUG - unexpected type value")
	}
}
