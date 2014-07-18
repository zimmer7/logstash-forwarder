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
	"flag"
	"fmt"
	"github.com/elasticsearch/kriterium/component/process"
	"github.com/elasticsearch/kriterium/errors"
	"github.com/elasticsearch/kriterium/flags"
	"github.com/elasticsearch/kriterium/panics"
	"log"
	"lsf"
)

// --------------------------------------------------------------
// shell bootup context
// --------------------------------------------------------------

// Binding names for command runner
var ContextKey = struct {
	WorkingDir, Debug, Verbose, Global, Environment string
}{
	WorkingDir:  "command.working-directory",
	Debug:       "command.flag.debug",
	Verbose:     "command.flag.verbose",
	Global:      "command.flag.global",
	Environment: "command.environment",
}

// --------------------------------------------------------------
// Command support types
// --------------------------------------------------------------

// Command Initializer func def
type InitFn func(env *lsf.Environment) error

// Command executor func def
type ExecFn func(env *lsf.Environment) error

// Command finalizer func def
type EndFn func(env *lsf.Environment) error

// --------------------------------------------------------------
// Command
// --------------------------------------------------------------

type Command struct {
	Name        string             // command name
	Flag        *flag.FlagSet      // command flag
	args        []string           // see ParseFlags()
	wd          string             // working directory - asserted non-empty
	EnvInit     bool               // flag: command initializes env
	Init        InitFn             // command init func - optional
	Run         ExecFn             // command run func
	End         EndFn              // command finalizer - optional
	IsActive    bool               // flag: command is active/long-running
	controller  process.Controller // active-command controller
	process     process.Process    // active-command process
	global      bool               // general option common to all commands
	debug       bool               // general option common to all commands
	verbose     bool               // general option common to all commands
	Option      interface{}        // command specific option (struct pointer)
	SubCommands []*Command         // sub-commands (if any - maybe nil)
}

// Initializes the command per bootup context and (cmd-line) args,
// sets command working directly, options, and flags.
//
// All args, parsed flags, and context bindings are asserted.
func (c *Command) Initialize(context map[string]interface{}, args ...string) (err error) {
	defer panics.Recover(&err)

	c.debug = context[ContextKey.Debug].(bool)
	c.verbose = context[ContextKey.Verbose].(bool)
	c.global = context[ContextKey.Global].(bool)

	// command line-options and args
	if c.Option != nil {
		c.Flag = flag.NewFlagSet(c.Name, flag.ContinueOnError)
		flags.MapStruct(c.Flag, c.Option)
		if e := c.ParseFlags(args...); e != nil {
			return e
		}
		c.args = c.Flag.Args()
	} else {
		c.args = args
	}

	wd := context[ContextKey.WorkingDir]
	if wd == nil {
		return errors.IllegalArgument("command.Initialize", "not bound", ContextKey.WorkingDir)
	}
	c.SetWorkingDirectory(wd.(string)) // panics on ""
	c.global = context[ContextKey.Global].(bool)

	return nil
}

// Used per --debug
func (c *Command) DebugOptions() {
	log.Printf("option: global? %t\n", c.global)
	log.Printf("option: verbose? %t\n", c.verbose)
	log.Printf("option: debug? %t\n", c.debug)
}

// (sub-)command args
func (c *Command) Args() []string {
	return c.args
}

// working directory
func (c *Command) Wd() string {
	return c.wd
}

// general command flag
func (c *Command) Global() bool {
	return c.global
}

// general command flag
func (c *Command) Debug() bool {
	return c.debug
}

// general command flag
func (c *Command) Verbose() bool {
	return c.verbose
}

// panics with errors.IllegalState if not an active command
func (c *Command) Controller() process.Controller {
	if !c.IsActive {
		panic(errors.IllegalState(c.Name, "is passive"))
	}
	return c.controller
}

// panics with errors.IllegalState if not an active command
func (c *Command) Process() process.Process {
	if !c.IsActive {
		panic(errors.IllegalState(c.Name, "is passive"))
	}
	return c.process
}

// Returns the named sub-command (if any). Internal use only.
//
// panics with errors.IllegalArgument if name does not match any sub-command.
func (c *Command) getSubCommand(name string) *Command {
	for _, cmd := range c.SubCommands {
		if cmd.Name == name {
			return cmd
		}
	}
	panic(errors.IllegalArgument("unknown-subcommand", name))
}

// panics will errors.IllegalArgument if working directory
// (arg wd) is blank
func (c *Command) SetWorkingDirectory(wd string) {
	if wd == "" {
		panic(errors.IllegalArgument("wd is zero-value"))
	}
	c.wd = wd
}

// panics with errors.Assertion if any of the required general
// properties of Command are not set.
func (c *Command) AssertSpec() {
	if c == nil {
		panic(errors.Assertion("command is nil"))
	}
	if c.Name == "" {
		panic(errors.Assertion("command.Name is not specified"))
	}
	if c.SubCommands == nil && c.Run == nil {
		panic(errors.Assertion("command.Run is nil (SubCommands is nil"))
	}
	if c.Option != nil && c.Flag == nil {
		panic(errors.Assertion("command.Flag is nil (Option is specified)"))
	}
}

// panics with errors.IllegalState if command is not property
// initialized. Successful return indicates command is ready for exec.
func (c *Command) AssertReadyState() {
	if c == nil {
		panic(errors.IllegalState("command is nil"))
	}
	if c.args == nil {
		panic(errors.IllegalState("command.arg is nil"))
	}
	//	if !c.EnvInit && c.env == nil {
	//		panic(errors.IllegalState("command.env is nil (command not an Env initilizer)"))
	//	}
	//	if c.Flag == nil {
	//		panic(errors.IllegalState("command.Flag is nil"))
	//	}
	if c.wd == "" {
		panic(errors.IllegalState("command.wd is zero-value"))
	}
}

// parses command flagset and sets command-line args (if any)
// returns errors.Usage on usage error
func (c *Command) ParseFlags(args ...string) error {
	if e := c.Flag.Parse(args); e != nil {
		return errors.Usage(c.Name)
	}
	//	c.args = c.Flag.Args()

	return nil
}

// --------------------------------------------------------------
// sub-commands
// --------------------------------------------------------------

// Executes the sub-command via the general command.Runner.
// Returned error is per command.Runner. See function for details.
func RunSubCommand(cmd, defaultSubCmd *Command, env *lsf.Environment) (err error) {
	defer panics.Recover(&err)

	args := cmd.Args()

	// generic TODO in command/command.go
	var subCmdName = defaultSubCmd.Name
	var subCmdArgs = []string{}
	var subCmdContext = cmd.SubCommandContext(env)
	if len(args) > 0 && MaybeSubCommand(args[0]) {
		subCmdName = args[0]
		subCmdArgs = args[1:]
	} else {
		subCmdArgs = args
	}

	var subCmd *Command = cmd.getSubCommand(subCmdName)

	return Runner(subCmd, subCmdContext, subCmdArgs...)
}

// Creates the command boot-up context for the child/sub-command
// per parent command's. Returned context (map) is non-nil and
// non-empty.
func (c *Command) SubCommandContext(env *lsf.Environment) map[string]interface{} {
	context := make(map[string]interface{})
	context[ContextKey.Environment] = env
	context[ContextKey.Debug] = c.debug
	context[ContextKey.Verbose] = c.verbose
	context[ContextKey.Global] = c.global
	context[ContextKey.WorkingDir] = c.wd
	return context
}

// Checks whether the arg may possibly represent the name of a sub-command.
func MaybeSubCommand(arg string) bool {
	return arg[0] != uint8('-')
}

// --------------------------------------------------------------
// Util
// --------------------------------------------------------------

// panics with errors.Assertion
func AssertStringProvided(info string, have, defval string) {
	if have == defval {
		panic(errors.Usage(fmt.Sprintf("option %q is required", info)))
	}
}

// panics with errors.Assertion
func AssertUint16Provided(info string, have, defval uint16) {
	if have == defval {
		panic(errors.Usage(fmt.Sprintf("option %q is required", info)))
	}
}

// Used per --debug. Emits via fmt the generic command properties.
func (c *Command) DebugCommand() {
	fmt.Printf("Command:\n")
	fmt.Printf("\tName:       %s\n", c.Name)
	fmt.Printf("\tWd:         %s\n", c.wd)
	fmt.Printf("\tglobal:     %t\n", c.global)
	fmt.Printf("\tdebug:      %t\n", c.debug)
	fmt.Printf("\tverbose:    %t\n", c.verbose)
}
