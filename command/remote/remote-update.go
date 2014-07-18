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

package remote

import (
	"fmt"
	"github.com/elasticsearch/kriterium/panics"
	"log"
	"lsf"
	"lsf/command"
	"lsf/schema"
)

var update *command.Command

func init() {
	update = &command.Command{
		Name:   "update",
		Run:    runUpdate,
		Option: &option,
	}
}

func runUpdate(env *lsf.Environment) (err error) {
	panics.Recover(&err)

	command.AssertStringProvided("remote-id", option.Id, "")

	updates := make(map[string][]byte)
	if option.Id != "" {
		updates[schema.PortElem.Id] = []byte(option.Id)
	}
	if option.Host != "" {
		updates[schema.PortElem.Host] = []byte(option.Host)
	}
	if option.Port != uint16(0) {
		updates[schema.PortElem.PortNum] = []byte(fmt.Sprintf("%d", option.Port))
	}

	if e := env.UpdateRemotePort(option.Id, updates); e != nil {
		return e
	}
	if update.Verbose() {
		log.Printf("Updated remote portal %q\n", option.Id)
	}
	return
}
