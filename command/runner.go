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
	"github.com/elasticsearch/kriterium/component/process"
	"github.com/elasticsearch/kriterium/errors"
	"github.com/elasticsearch/kriterium/panics"
	"log"
	"lsf"
	"os"
)

// use is limited to this file only - convenience
var verbose bool

// use is limited to this file only - convenience
var debug bool

// Command runner. Initializes and executes both active
// and passive commands.
//
// Input args must not be nil.
//
func Runner(command *Command, context map[string]interface{}, args ...string) (err error) {
	defer panics.Recover(&err)

	// be verbose if explicitly required or if in debug mode
	// debug may emit additional outputs during command run.
	debug = context[ContextKey.Debug].(bool)
	verbose = debug || context[ContextKey.Verbose].(bool)

	/// initialize //////////////////////////////////////////////

	if context == nil {
		return errors.IllegalArgument("command.Runner", "context is nil")
	}
	if e := command.Initialize(context, args...); e != nil {
		return e
	}

	command.AssertSpec() // panics

	/// get ready ///////////////////////////////////////////////

	// REVU seems all this really belongs to command/command.go
	//      and command initialization
	// TODO begin

	// create command-control if command is active.
	// spec command.Stop on signals os.Kill and os.Interrupt
	if command.IsActive {
		cmdCnc := process.NewCmdCtl()
		cmdCnc.CommandOnSignal(process.Stop, os.Interrupt, os.Kill)
		command.controller = cmdCnc.Controller()
		command.process = cmdCnc.Process()
	}

	// create and set command environment
	// If command is itself the env initializer, skip it.
	var env *lsf.Environment
	if v := context[ContextKey.Environment]; v != nil {
		env = v.(*lsf.Environment)
	} else if !command.EnvInit {
		env = lsf.NewEnvironment()
		e := env.Initialize(command.Wd())
		panics.OnError(e, "command.Run:", "env.Initialize:")

		// env is an active component.
		// shut it down when runner returns.
		// REVU: TODO: this should use kriterium process Controller ..
		//       TODO: mod lsf.Environment per above.
		//       TODO: defer env.Process.Command(process.Stop)
		defer func() {
			env.Shutdown()
		}()
	}

	// TODO: END

	command.AssertReadyState() // panics

	/// action //////////////////////////////////////////////////

	// run command initializer (if any)
	if command.Init != nil {
		e := command.Init(env)
		println(e)
		panics.OnError(e)
	}

	// run the command proper
	if e := command.Run(env); e != nil {
		return e
	}

	/// and cut /////////////////////////////////////////////////

	if command.End == nil {
		return
	}

	return command.End(env)
}

// emits per --verbose (and logical verbose) command shell settings.
// 'must' arg forces emit.
func emit(must bool, msg string) {
	if must || verbose {
		log.Println(msg)
	}
}
