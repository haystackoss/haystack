//go:build darwin
// +build darwin

package limit

import (
	"fmt"
	"os"
	"os/exec"
)

func run(args []string) ([]byte, int, error) {
	cmd := exec.Command(args[0], args[1:]...)
	stdout, err := cmd.CombinedOutput()
	exitCode := cmd.ProcessState.ExitCode()
	return stdout, exitCode, err
}

const (
	maxFilesPerProc = 99999
	maxFiles        = 200000
)

func InitLimit() {
	//sudo sysctl -w kern.maxfiles=49152
	// sudo sysctl -w kern.maxfilesperproc=24576
	cmd0 := []string{"sudo", "-n", "true", "2>/dev/null;"}
	cmd1 := []string{"sudo", "sysctl", "-w", "kern.maxfiles=200000"}
	cmd1 := []string{"sudo", "sysctl", "-w", "kern.maxfilesperproc=99999"}
	cmd4 := []string{"sudo", "sysctl", "-p"}
	_, exitCode, _ := run(cmd0)
	if exitCode != 0 {
		fmt.Printf("We would like to increase kern.maxfilesperproc to %d\n", maxFilesPerProc) // just a warning before the prompt
	}

	_, _, err1 := run(cmd1)
	_, _, err2 := run(cmd2)
	_, _, err3 := run(cmd3)
	_, _, err4 := run(cmd4)

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		fmt.Println("Failed to increase inotify limit")
		os.Exit(1)
	}

}
