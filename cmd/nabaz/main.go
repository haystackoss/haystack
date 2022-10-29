package main

import (
	"fmt"
	"io"
	_ "net/http/pprof"
	"os"

	"github.com/nabaz-io/argparse"
	"github.com/nabaz-io/nabaz/pkg/fixme"
	"github.com/nabaz-io/nabaz/pkg/testrunner"
	"github.com/nabaz-io/nabaz/pkg/version"
	"github.com/pkg/profile"
)

type Subcommand int

const (
	TEST Subcommand = iota
	FIXME
	VERSION
	UNKNOWN
)

type ProgramArguments struct {
	Test  testrunner.Arguements
	FixMe fixme.Arguements
	cmd   Subcommand
}

func ParseArguements(args []string) *ProgramArguments {
	command := UNKNOWN
	cliParser := argparse.NewParser("nabaz", "")
	testCmd := cliParser.NewCommand("test", "Runs tests")
	fixmeCmd := cliParser.NewCommand("fixme", "A fixme list of broken tests")
	versionCmd := cliParser.NewCommand("version", "Gets version of nabaz.")

	cmdline := fixmeCmd.String("", "cmdline", &argparse.Options{
		Required: true,
		Help:     "i.e: go test ./...",
	})
	cmdlineTestRunner := testCmd.String("", "cmdline", &argparse.Options{
		Required: true,
		Help:     "i.e: go test ./...",
	})
	pkgs := testCmd.String("", "pkgs", &argparse.Options{
		Required: false,
		Help:     "list packages being tested in go tested",
	})
	// Naming a positional arguement isn't
	repoPath := fixmeCmd.StringPositional("repo_path", &argparse.Options{
		Required: false,
		Help:     "Postional arguement (don't use flag)",
		Default:  ".",
	})
	// Naming a positional arguement isn't
	repoPathTestRunner := testCmd.StringPositional("repo_path", &argparse.Options{
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
		Test: testrunner.Arguements{
			Cmdline:  *cmdlineTestRunner,
			Pkgs:     *pkgs,
			RepoPath: *repoPathTestRunner,
		},
		FixMe: fixme.Arguements{
			Cmdline:  *cmdline,
			RepoPath: *repoPath,
		},

		cmd: command,
	}
}

func NotImplemented() {
	fmt.Println("Not Implemented")
	os.Exit(1)
}

func isDebug() bool {
	val := os.Getenv("DEBUG_NABAZ")
	for _, yesVal := range []string{"1", "true", "yes"} {
		if val == yesVal {
			return true
		}
	}
	return false
}

func run(args []string, stdout io.Writer) {
	if isDebug() {
		defer profile.Start(profile.CPUProfile).Stop()
	}

	parsedArgs := ParseArguements(args)
	switch parsedArgs.cmd {
	case TEST:
		testRunnerArgs := &parsedArgs.Test
		testrunner.Execute(testRunnerArgs)
	case FIXME:
		fixmeArgs := &parsedArgs.FixMe
		fixme.Execute(fixmeArgs)
	case VERSION:
		version.Execute()
	default:
		fmt.Println("Unknown command")
	}
}

func main() {
	run(os.Args, os.Stdout)
}
