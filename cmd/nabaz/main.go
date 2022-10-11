package main

import (
	"fmt"
	"io"
	"os"

	"github.com/nabaz-io/argparse"
	"github.com/nabaz-io/nabaz/pkg/testrunner"
	"github.com/nabaz-io/nabaz/pkg/version"
)

type Subcommand int

const (
	TEST Subcommand = iota
	VERSION
	UNKNOWN
)

type ProgramArguments struct {
	Test testrunner.Arguements
	cmd  Subcommand
}

func ParseArguements(args []string) *ProgramArguments {
	command := UNKNOWN
	cliParser := argparse.NewParser("nabaz", "")
	testCmd := cliParser.NewCommand("test", "Runs tests")
	versionCmd := cliParser.NewCommand("version", "Gets version of nabaz.")

	cmdline := testCmd.String("", "cmdline", &argparse.Options{
		Required: true,
		Help:     "i.e: go test ./...",
	})
	pkgs := testCmd.String("", "pkgs", &argparse.Options{
		Required: false,
		Help:     "list packages being tested in go tested",
	})
	// Naming a positional arguement isn't
	repoPath := testCmd.StringPositional("repo_path", &argparse.Options{
		Required: false,
		Help:     "Postional arguement (don't use flag)",
		Default:  ".",
	})

	err := cliParser.Parse(args)

	if err != nil {
		fmt.Print(testCmd.Usage(err))
		os.Exit(1)
	}

	switch {
	case versionCmd.Happened():
		command = VERSION
	case testCmd.Happened():
		command = TEST
	default:
		command = UNKNOWN
	}

	return &ProgramArguments{
		Test: testrunner.Arguements{
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
	case TEST:
		testRunnerArgs := &parsedArgs.Test
		testrunner.Execute(testRunnerArgs)
	case VERSION:
		version.Execute()
	default:
		fmt.Println("Unknown command")
	}
}

func main() {
	run(os.Args, os.Stdout)
}
