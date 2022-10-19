package main

import (
	"fmt"
	"io"
	"os"

	"github.com/nabaz-io/argparse"
	"github.com/nabaz-io/nabaz/pkg/fixme"
	"github.com/nabaz-io/nabaz/pkg/version"
)

type Subcommand int

const (
	TEST  Subcommand = iota
	FIXME Subcommand = iota
	VERSION
	UNKNOWN
)

type ProgramArguments struct {
	Test fixme.Arguements
	cmd  Subcommand
}

func ParseArguements(args []string) *ProgramArguments {
	command := UNKNOWN
	cliParser := argparse.NewParser("nabaz", "")
	testCmd := cliParser.NewCommand("test", "Runs tests")
	versionCmd := cliParser.NewCommand("version", "Gets version of nabaz.")
	fixmeCmd := cliParser.NewCommand("fixme", "Fixme list of tests.")

	cmdline := fixmeCmd.String("", "cmdline", &argparse.Options{
		Required: true,
		Help:     "i.e: go test ./...",
	})
	pkgs := fixmeCmd.String("", "pkgs", &argparse.Options{
		Required: false,
		Help:     "list packages being tested in go tested",
	})
	// Naming a positional arguement isn't
	repoPath := fixmeCmd.StringPositional("repo_path", &argparse.Options{
		Required: false,
		Help:     "Postional arguement (don't use flag)",
		Default:  ".",
	})

	err := cliParser.Parse(args)

	if err != nil {
		fmt.Print(fixmeCmd.Usage(err))
		os.Exit(1)
	}

	switch {
	case versionCmd.Happened():
		command = VERSION
	case testCmd.Happened():
		command = TEST
	case fixmeCmd.Happened():
		command = FIXME
	default:
		command = UNKNOWN
	}

	return &ProgramArguments{
		Test: fixme.Arguements{
			Cmdline:  *cmdline,
			Pkgs:     *pkgs,
			RepoPath: *repoPath,
		},
		cmd: command,
	}
}

func NotImplemented() {
	fmt.Println("Not Implemented")
	os.Exit(1)
}

func run(args []string, stdout io.Writer) {

	parsedArgs := ParseArguements(args)
	switch parsedArgs.cmd {
	case FIXME:
		fixmeArgs := &parsedArgs.Test
		fixme.Execute(fixmeArgs)
	case TEST:
		panic("not supported")
	case VERSION:
		version.Execute()
	default:
		fmt.Println("Unknown command")
	}
}

func main() {
	run(os.Args, os.Stdout)
}
