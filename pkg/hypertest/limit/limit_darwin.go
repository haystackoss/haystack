//go:build darwin

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
	cmd0 := []string{"sudo", "-n", "true", "2>/dev/null;"}
	cmd1 := []string{"sudo", "sysctl", "-w", "kern.maxfiles=200000"}
	cmd2 := []string{"sudo", "sysctl", "-w", "kern.maxfilesperproc=99999"}
	_, exitCode, _ := run(cmd0)
	if exitCode != 0 {
		fmt.Println("Let's increase kern.maxfilesperproc")
	}

	_, _, err1 := run(cmd1)
	_, _, err2 := run(cmd2)

	if err1 != nil || err2 != nil {
		fmt.Println("Failed to increase inotify limit")
		os.Exit(1)
	}

}
