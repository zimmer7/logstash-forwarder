package main

import (
	"flag"
	"github.com/elasticsearch/kriterium/errors"
	"github.com/elasticsearch/kriterium/flags"
	"github.com/elasticsearch/kriterium/panics"
	"log"
	"lsf/command"
	"lsf/command/give"
	"lsf/command/help"
	"lsf/command/initialize"
	"lsf/command/migrate"
	"lsf/command/remote"
	"lsf/command/stream"
	"lsf/command/track"
	"os"
)

// REVU: this really belongs to command/command.go
var commands = []*command.Command{
	initialize.Command,
	stream.Command,
	remote.Command,
	track.Command,
	give.Command,
	help.Command,
	migrate.Command,
}

const (
	stat_ok int = iota
	stat_usage
	stat_error
)

var option = struct {
	Verbose bool   `long:verbose about:"verbose ops and logs"`
	Global  bool   `long:global  about:"apply command globally"`
	Debug   bool   `long:debug   about:"debug ops and logs"`
	Home    string `long:home    about:"set working directory"`
	About   bool   `long:about   about:"display app info and exit"`
	Version bool   `long:version about:"display version and exit"`
}{}

func init() {
	log.SetFlags(0)

	wd, e := os.Getwd()
	panics.OnError(e)

	option.Home = wd

	e = flags.MapStruct(flag.CommandLine, &option)
	if e != nil {
		log.Fatal(e)
	}
}

var verbose bool

func main() {
	flag.Parse()
	verbose = option.Debug || option.Verbose

	if option.Debug {
		panics.DEBUG = true
	}

	debugOptions()

	switch {
	case option.About:
		exit(onAbout())
	case option.Version:
		exit(onVersion())
	}

	args := flag.Args()
	if len(args) < 1 {
		exit(onUsage())
	}

	cmd := getCommandByName(args[0])
	if cmd == nil {
		exit(onUsage())
	}

	context := make(map[string]interface{})
	context[command.ContextKey.Debug] = option.Debug
	context[command.ContextKey.Verbose] = option.Verbose
	context[command.ContextKey.Global] = option.Global
	context[command.ContextKey.WorkingDir] = option.Home

	stat := 0
	if e := command.Runner(cmd, context, args[1:]...); e != nil {
		switch {
		case errors.Usage.Matches(e):
			emit(true, e.Error())
			stat = stat_usage
		default:
			stat = onError(e)
		}
	}

	exit(stat)
}

func exit(stat int) {
	os.Exit(stat)
}

func emit(must bool, msg string) {
	if must || verbose {
		log.Println(msg)
	}
}

func onError(e error) int {
	emit(true, e.Error())
	return stat_error
}

func onUsage() int {
	emit(true, "usage: lsf [<options>] [<command> <options>]")
	return stat_usage
}

func onVersion() int {
	emit(true, "todo: version msg")
	return stat_ok
}

func onAbout() int {
	emit(true, "todo: about msg")
	return stat_ok
}

func getCommandByName(name string) *command.Command {
	for _, cmd := range commands {
		if cmd.Name == name {
			return cmd
		}
	}
	return nil
}

/// temp debug /////////////////////////////

func debugCmdLineArgs(args ...string) {
	for _, arg := range args {
		log.Printf("arg: %s\n", arg)
	}
}
func debugOptions() {
	if !option.Debug {
		return
	}

	log.Printf("verbose? %t\n", option.Verbose)
	log.Printf("global?  %t\n", option.Global)
	log.Printf("version? %t\n", option.Version)
	log.Printf("debug?   %t\n", option.Debug)
	log.Printf("about?   %t\n", option.About)
	log.Printf("home?    %s\n", option.Home)
}
