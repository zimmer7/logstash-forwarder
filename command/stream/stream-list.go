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
	"github.com/elasticsearch/kriterium/panics"
	"log"
	"lsf"
	"lsf/command"
	"lsf/schema"
)

var list *command.Command
var listOption = struct {
	Id   string `code:s long:stream-id about:"stream id"`
	Info bool   `code:i long:info about:"show full stream info"`
}{}

func init() {
	list = &command.Command{
		Name:   "list",
		Run:    runList,
		Option: &listOption,
	}
}

func runList(env *lsf.Environment) (err error) {
	defer panics.Recover(&err)

	digests := env.GetResourceDigests("stream", listOption.Info, schema.LogStreamDigest)
	for _, digest := range digests {
		log.Println(digest)
	}

	if list.Verbose() && len(digests) == 0 {
		log.Printf("There are no streams defined in %s\n", list.Wd())
	}
	return
}
