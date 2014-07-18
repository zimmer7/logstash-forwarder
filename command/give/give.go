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

package give

import (
	"github.com/elasticsearch/kriterium/panics"
	"log"
	"lsf"
	"lsf/command"
)

var Command *command.Command

var option = struct {
	N uint `long:me about:"Ask and you shall receive!"`
}{
	N: 3,
}

func init() {
	Command = &command.Command{
		Name:   "give",
		Option: &option,
		Run:    run,
	}
}

func run(env *lsf.Environment) (err error) {
	defer panics.Recover(&err)

	for n := 0; n < int(option.N); n++ {
		log.Printf("h%c gs\n", rune('\u2661'))
	}
	return
}
