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
	"github.com/elasticsearch/kriterium/panics"
	"lsf"
	"lsf/schema"
)

const addStreamCmdCode lsf.CommandCode = "stream-add"

var addStream *lsf.Command
var addStreamOptions *editStreamOptions

func init() {

	flagset := FlagSet(addStreamCmdCode)
	addStream = &lsf.Command{
		Name:  addStreamCmdCode,
		About: "Add a new log stream",
		Init:  initAddStream,
		Run:   runAddStream,
		Flag:  flagset,
	}
	addStreamOptions = &editStreamOptions{
		Verbose: flags.NewBoolOption(flagset, "v", "verbose", false, "be verbose in list", false),
		Global:  flags.NewBoolOption(flagset, "G", "global", false, "global scope flag for command", false),
		Id:      flags.NewStringOption(flagset, "s", "stream-id", "", "unique identifier for stream", true),
		Path:    flags.NewStringOption(flagset, "p", "path", "", "path to log files", true),
		Mode:    flags.NewStringOption(flagset, "m", "journal-mode", "", "stream journaling mode (rotation|rollover)", true),
		Pattern: flags.NewStringOption(flagset, "n", "name-pattern", "", "naming pattern of journaled log files", true),
	}
}

func initAddStream(env *lsf.Environment, args ...string) (err error) {
	e := flags.UsageVerify(addStreamOptions)
	if e != nil {
		return e
	}

	mode := addStreamOptions.Mode.Get()
	switch schema.ToJournalModel(mode) {
	case schema.JournalModel.Rotation, schema.JournalModel.Rollover:
	default:
		println("HERE")
		return errors.Usage("stream-add", "option", "option mode must be one of {rollover, rotation}")
	}
	return
}

func runAddStream(env *lsf.Environment, args ...string) (err error) {
	defer panics.Recover(&err)

	e := flags.UsageVerify(addStreamOptions)
	panics.OnError(e)

	id := addStreamOptions.Id.Get()
	pattern := addStreamOptions.Pattern.Get()
	journalMode := addStreamOptions.Mode.Get()
	basepath := addStreamOptions.Path.Get()
	fields := make(map[string]string) // TODO: fields needs a solution

	return env.AddLogStream(id, basepath, pattern, journalMode, fields)
}
