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
	"github.com/elasticsearch/kriterium/panics"
	"lsf"
	"lsf/schema"
)

const updateStreamCmdCode lsf.CommandCode = "stream-update"

var updateStream *lsf.Command
var updateStreamOptions *editStreamOptions

func init() {

	flagset := FlagSet(updateStreamCmdCode)
	updateStream = &lsf.Command{
		Name:  updateStreamCmdCode,
		About: "Update a new log stream",
		Init:  initUpdateStream,
		Run:   runUpdateStream,
		Flag:  flagset,
	}
	updateStreamOptions = &editStreamOptions{
		Verbose: flags.NewBoolOption(flagset, "v", "verbose", false, "be verbose in list", false),
		Global:  flags.NewBoolOption(flagset, "G", "global", false, "global scope flag for command", false),
		Id:      flags.NewStringOption(flagset, "s", "stream-id", "", "unique identifier for stream", true),
		Path:    flags.NewStringOption(flagset, "p", "path", "", "path to log files", false),
		Mode:    flags.NewStringOption(flagset, "m", "journal-mode", "", "stream journaling mode (rotation|rollover)", false),
		Pattern: flags.NewStringOption(flagset, "n", "name-pattern", "", "naming pattern of journaled log files", false),
	}
}

func initUpdateStream(env *lsf.Environment, args ...string) (err error) {
	return flags.UsageVerify(updateStreamOptions)
}

func runUpdateStream(env *lsf.Environment, args ...string) (err error) {
	defer panics.Recover(&err)

	id := updateStreamOptions.Id.Get()
	updates := make(map[string][]byte)

	// update stream config document
	var option flags.StringOption
	option = *updateStreamOptions.Pattern
	if option.Provided() {
		v := []byte(option.Get())
		updates[schema.LogStreamElem.Pattern] = v
	}
	option = *updateStreamOptions.Path
	if option.Provided() {
		v := []byte(option.Get())
		updates[schema.LogStreamElem.BasePath] = v
	}
	option = *updateStreamOptions.Mode
	if option.Provided() {
		v := []byte(option.Get())
		updates[schema.LogStreamElem.JournalModel] = v
	}

	return env.UpdateLogStream(id, updates)

}
