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
	"log"
	"lsf"
	"lsf/command"
)

var Command *command.Command

// common to remote -update -add.
var option = struct {
	Id   string `code:r long:remote-id   about:"remote portal unique identifier"`
	Host string `code:h long:host-name   about:"remote portal host's name"`
	Port uint16 `code:p long:port-number about:"remote portal host's port number"`
}{}

func init() {
	Command = &command.Command{
		Name:        "remote",
		Run:         run,
		SubCommands: []*command.Command{list, add, remove, update},
	}
}

func run(env *lsf.Environment) (err error) {

	//	Command.DebugCommand()
	return command.RunSubCommand(Command, list, env)
}

// temp debug remote.option
func debugOptions() {
	log.Printf("id?      %s\n", option.Id)
	log.Printf("host?    %s\n", option.Host)
	log.Printf("portnum? %s\n", option.Port)
}
