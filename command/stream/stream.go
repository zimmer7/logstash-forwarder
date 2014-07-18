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

package stream

import (
	"log"
	"lsf"
	"lsf/command"
)

var Command *command.Command

// common to stream -update -add.
var option = struct {
	Id      string `code:s long:stream-id about:"stream id"`
	Path    string `code:p long:path about:"log files' basepath"`
	Pattern string `code:n long:name-pattern about:"log files' naming pattern"`
	Model   string `code:m long:rotation-model about:"log journaling mode (rotation|rollover)"`
}{}

func init() {
	Command = &command.Command{
		Name:        "stream",
		Run:         run,
		SubCommands: []*command.Command{list, add, remove, update},
	}
}

func run(env *lsf.Environment) (err error) {

	//	Command.DebugCommand()
	return command.RunSubCommand(Command, list, env)
}

// temp debug stream.option
func debugOptions() {
	log.Printf("id?      %s\n", option.Id)
	log.Printf("path?    %s\n", option.Path)
	log.Printf("pattern? %s\n", option.Pattern)
	log.Printf("model?   %s\n", option.Model)
	//	log.Printf("verbose? %t\n", option.Verbose)

}
