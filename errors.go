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

package lsf

import (
	"github.com/elasticsearch/kriterium/errors"
)

// ----------------------------------------------------------------------------
// error codes
// ----------------------------------------------------------------------------

// REVU: All these can be either generic IllegalState or Fatal error types.
// TODO: remove these
var LsfError = struct {
	OpFailure,
	EnvironmentExists,
	EnvironmentDoesNotExist,
	ResourceExists,
	ResourceDoesNotExist,
	_stub errors.TypedError
}{
	OpFailure:               errors.New("lsf operation failed"),
	EnvironmentExists:       errors.New("lsf environment already exists"),
	EnvironmentDoesNotExist: errors.New("lsf environment does not exists at location"),
	ResourceExists:          errors.New("lsf resource already exists"),
	ResourceDoesNotExist:    errors.New("lsf resource does not exist"),
}

var WARN = struct {
	NoOp errors.TypedError
}{
	NoOp: errors.New("warning: no op"),
}
