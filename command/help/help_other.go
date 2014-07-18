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

package help

import (
	"fmt"
	"github.com/elasticsearch/kriterium/errors"
	"github.com/elasticsearch/kriterium/panics"
	"lsf/command"
	"os"
	"os/exec"
)

func help() (err error) {

	defer panics.Recover(&err)

	// REVU: the man pages should be included in the binary
	// TODO: the first run of lsf should install the man pages
	//       in the global lsf portal. TODO in command/initialization/
	//
	//	manpath := lsf.AbsolutePath(Command.Wd())
	subject := "git" // for now .. TODO replace with lsf
	args := Command.Args()
	if len(args) > 0 && command.MaybeSubCommand(args[0]) {
		subject = fmt.Sprintf("%s-%s", subject, args[0])
	}

	var synopsis = ""
	if option.WhatIs {
		synopsis = "-f"
	}
	man := exec.Command("man", synopsis, subject) // TODO: add -M manpath
	man.Stdout = os.Stdout
	man.Stderr = os.Stderr

	if e := man.Start(); e != nil {
		return errors.Fatal(e)
	}
	if e := man.Wait(); e != nil {
		return errors.Fatal(e)
	}

	return
}
