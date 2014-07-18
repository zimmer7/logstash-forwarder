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

package initialize

import (
	"github.com/elasticsearch/kriterium/panics"
	"log"
	"lsf"
	"lsf/command"
)

var Command *command.Command

var option = struct {
	Force bool `code:f long:force about:"force operation"`
}{}

func init() {

	Command = &command.Command{
		Name:    "init",
		Option:  &option,
		EnvInit: true,
		Run:     run,
	}
}

func run(env *lsf.Environment) (err error) {
	defer panics.Recover(&err)

	home, e := lsf.AbsolutePath(Command.Wd())
	if e != nil {
		return e
	}

	action := "Initialized"
	if env.Exists(home) {
		panics.OnFalse(option.Force, "lsf init:", "existing environment. use -force flag to reinitialize")
		action = "Re-Initialize"
	}

	portalPath, e := lsf.CreateEnvironment(home, option.Force)
	panics.OnError(e, "lsf init:")

	if Command.Verbose() {
		log.Printf("%s LS/F portal at %s\n", action, portalPath)
	}

	return
}
