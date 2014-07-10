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
	"github.com/elasticsearch/kriterium/flags"
	"lsf"
)

const cmd_stream lsf.CommandCode = "stream"

type editStreamOptions struct {
	Verbose, Global         *flags.BoolOption
	Id, Mode, Path, Pattern *flags.StringOption
}

var Stream *lsf.Command
var streamOption *flags.BoolOption

const (
	streamOptionVerbose   = "command.stream.option.verbose"
	streamOptionGlobal    = "command.stream.option.global"
	streamOptionsSelected = "command.stream.option.selected"
)

func init() {

	flagset := FlagSet(cmd_stream)
	Stream = &lsf.Command{
		Name:  cmd_stream,
		About: "Stream is a top level command for log stream configuration and management",
		Run:   runStream,
		Flag:  flagset,
	}
	streamOption = flags.NewBoolOption(flagset, "v", "verbose", false, "be verbose in list", false)
}

func runStream(env *lsf.Environment, args ...string) error {

	env.Set(streamOptionVerbose, streamOption.Get())

	xoff := 0
	var subcmd *lsf.Command = listStream
	if len(args) > 0 {
		subcmd = getSubCommand(args[0])
		xoff = 1
	}

	return lsf.Run(env, subcmd, args[xoff:]...)
}

func getSubCommand(subcmd string) *lsf.Command {

	var cmd *lsf.Command
	switch lsf.CommandCode("stream-" + subcmd) {
	case addStreamCmdCode:
		cmd = addStream
	case removeStreamCmdCode:
		cmd = removeStream
	case updateStreamCmdCode:
		cmd = updateStream
	case listStreamCmdCode:
		cmd = listStream
	default:
		// not panic -- return error TODO
		panic("BUG - unknown subcommand for stream: " + subcmd)
	}
	return cmd
}
