//go:build linux

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
	maxUserWatches   = 124983
	maxUserInstances = 10000
)

func InitLimit() {
	cmd1 := []string{"sudo", "sysctl", fmt.Sprintf("fs.inotify.max_user_watches=%d", maxUserWatches)}
	cmd2 := []string{"sudo", "sysctl", fmt.Sprintf("fs.inotify.max_user_instances=%d", maxUserInstances)}
	cmd3 := []string{"sudo", "sysctl", "-p"}

	cmd0 := []string{"sudo", "-n", "true", "2>/dev/null;"}
	_, exitCode, _ := run(cmd0)
	if exitCode != 0 {
		fmt.Printf("We would like to increase fs.inotify limits to %d\n", maxUserWatches) // just a warning before the prompt
	}

	_, _, err1 := run(cmd1)
	_, _, err2 := run(cmd2)
	_, _, err3 := run(cmd3)

	if err1 != nil || err2 != nil || err3 != nil {
		fmt.Println("Failed to increase inotify limit")
		os.Exit(1)
	}

}
