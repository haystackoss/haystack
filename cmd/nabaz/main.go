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
	storageUrl := testCmd.String("", "storage-url", &argparse.Options{
		Required: true,
		Help:     "The nabaz storage. [required]",
	})
	webUrl := testCmd.String("", "web-url", &argparse.Options{
		Required: true,
		Help:     "Url to nabaz web-console. [required]",
	})
	pkgs := testCmd.String("", "pkgs", &argparse.Options{
		Required: false,
		Help:     "list packages being tested in go tested",
	})
	token := testCmd.String("", "github-token", &argparse.Options{
		Required: false,
		Help:     "Token to access github account.",
	})
	repoUrl := testCmd.String("", "repo-url", &argparse.Options{
		Required: false,
		Help:     "The repo url to work with, i.e https://github.com/trovalds/linux",
	})
	commitID := testCmd.String("", "commit-id", &argparse.Options{
		Required: false,
		Help:     "The commit id of the change (in case there is no .git)",
	})

	username := testCmd.String("", "username", &argparse.Options{
		Required: false,
		Help:     "Username for storage.",
	})

	password := testCmd.String("", "password", &argparse.Options{
		Required: false,
		Help:     "Password for storage.",
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
			Cmdline:    *cmdline,
			StorageUrl: *storageUrl,
			WebUrl:     *webUrl,
			Pkgs:       *pkgs,
			Token:      *token,
			RepoUrl:    *repoUrl,
			CommitID:   *commitID,
			Username:   *username,
			Password:   *password,
			RepoPath:   *repoPath,
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
