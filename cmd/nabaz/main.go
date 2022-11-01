package main

import (
	"fmt"
	"io"
	_ "net/http/pprof"
	"os"

	"github.com/nabaz-io/argparse"
	"github.com/nabaz-io/nabaz/pkg/hypertest"
	"github.com/nabaz-io/nabaz/pkg/testrunner"
	"github.com/nabaz-io/nabaz/pkg/version"
	"github.com/pkg/profile"
)

type Subcommand int

const (
	TEST Subcommand = iota
	HYPERTEST
	VERSION
	UNKNOWN
)

type ProgramArguments struct {
	Test      testrunner.Arguements
	HyperTest hypertest.Arguements
	cmd       Subcommand
}

func ParseArguements(args []string) *ProgramArguments {
	command := UNKNOWN
	cliParser := argparse.NewParser("nabaz", "")
	testCmd := cliParser.NewCommand("test", "Runs tests")
	hypertestCmd := cliParser.NewCommand("hypertest", "Generates an auto-updating todo list of coding subtasks.")
	versionCmd := cliParser.NewCommand("version", "Gets version of nabaz.")

	cmdline := hypertestCmd.String("", "cmdline", &argparse.Options{
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
	repoPath := hypertestCmd.StringPositional("repo_path", &argparse.Options{
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
		fmt.Print(hypertestCmd.Usage(err))
		os.Exit(1)
	}

	switch {
	case versionCmd.Happened():
		command = VERSION
	case testCmd.Happened():
		command = TEST
	case hypertestCmd.Happened():
		command = HYPERTEST
	default:
		command = UNKNOWN
	}

	return &ProgramArguments{
		Test: testrunner.Arguements{
			Cmdline:  *cmdlineTestRunner,
			Pkgs:     *pkgs,
			RepoPath: *repoPathTestRunner,
		},
		HyperTest: hypertest.Arguements{
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
	case HYPERTEST:
		hypertestArgs := &parsedArgs.HyperTest
		hypertest.Execute(hypertestArgs)
	case VERSION:
		version.Execute()
	default:
		fmt.Println("Unknown command")
	}
}

func main() {
	run(os.Args, os.Stdout)
}
