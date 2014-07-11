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

package command

import (
	"github.com/elasticsearch/kriterium/errors"
	"github.com/elasticsearch/kriterium/flags"
	"lsf"
)

const cmd_debug lsf.CommandCode = "debug"

var Debug *lsf.Command
var debugOption *flags.StringOption

func init() {

	flagset := FlagSet(cmd_debug)
	Debug = &lsf.Command{
		Name:  cmd_debug,
		About: "Provides usage information for LS/F commands",
		Init:  initDebug,
		Run:   runDebug,
		Flag:  flagset,
		Usage: "debug <command>",
	}
	debugOption = flags.NewStringOption(flagset, "c", "command", "", "the command to debug", true)
}

func initDebug(env *lsf.Environment, args ...string) error {
	return flags.VerifyRequiredOption(debugOption)
}

func runDebug(env *lsf.Environment, args ...string) error {
	//	_ = debugOption.Get()
	return errors.NotImplemented("command:", cmd_debug)
}
