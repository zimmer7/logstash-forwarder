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

package track

import (
	"github.com/elasticsearch/kriterium/component/process"
	"github.com/elasticsearch/kriterium/errors"
	"github.com/elasticsearch/kriterium/panics"
	"lsf"
	"lsf/command"
	"lsf/fs"
	"lsf/lsproc"
	"time"
)

var Command *command.Command

var option = struct {
	StreamId string        `code:s long:stream-id    about:"stream id"`
	Delay    time.Duration `code:d long:update-delay about:"delay between track events"`
	MaxSize  uint16        `code:n long:max-size     about:"max size of the track FS Object state cache"`
	MaxAge   fs.InfoAge    `code:a long:max-age      about:"max age of item in the track FS Object state cache"`
}{
	Delay: time.Second, // 1 snapshot per second
}

func init() {
	Command = &command.Command{
		Name:     "track",
		Option:   &option,
		Run:      run,
		IsActive: true,
	}
}

func run(env *lsf.Environment) (err error) {
	panics.Recover(&err)

	// TODO: check that command.go checks controller/process in ReadyState()
	command.AssertStringProvided("stream-id", option.StreamId, "")

	//	trackConfig := &lsproc.TrackConfig{env, Command.Debug(), Command.Verbose(), option.StreamId, option.Delay, option.MaxSize, option.MaxAge}
	trackConfig := lsproc.NewDefaultConfig(env, option.StreamId)
	trackConfig.Debug = Command.Debug()
	trackConfig.Verbose = Command.Verbose()

	go lsproc.TrackProcess(Command.Controller(), trackConfig)

	// REVU: TODO: (low priority) refactor to kriterium -- BEGIN

	var cmdResponse interface{}

	// start tracking process
	Command.Process().Signal() <- process.Start
	cmdResponse = <-Command.Process().Response()
	switch t := cmdResponse.(type) {
	case process.CommandCode:
		if t != process.Start {
			return errors.IllegalState("unexpected response from process:", cmdResponse)
		}
	case error:
		return errors.Fatal(t.Error())
	}

	// wait until user signals end by os.Kill or os.Interrupt
	cmdResponse = <-Command.Process().Response()
	switch t := cmdResponse.(type) {
	case process.CommandCode:
		if t != process.Stop {
			return errors.IllegalState("unexpected response from process:", cmdResponse)
		}
	case error:
		return errors.Fatal(t.Error())
	}
	// REVU: TODO: (low priority) refactor to kriterium -- END

	return
}
