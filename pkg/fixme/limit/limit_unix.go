//go:build !windows
// +build !windows

package limit

import (
	"fmt"
	"os"
	"syscall"
)

func InitLimit() {
	var rLimit syscall.Rlimit
	rLimit.Max = 124983
	rLimit.Cur = 124983
	err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		fmt.Println("Error Setting Rlimit ", err)
		os.Exit(1)
	}
	
	// cmd1 := []string{"sudo", "sysctl", "fs.inotify.max_user_watches=124983"}
    // cmd2 := []string{"sudo", "sysctl", "fs.inotify.max_user_instances=10000"}
	// cmd3 := []string{"sudo", "sysctl", "fs.inotify.max_queued_events=10000"}
	// cmd4 := []string{"sudo", "sysctl", "-p"}
	// _, _, err1 := run(cmd1)
	// _, _, err2 := run(cmd2)
	// _, _, err3 := run(cmd3)
	// _, _, err4 := run(cmd4)

	// if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
	// 	fmt.Println("Failed to increase inotify limit")
	// 	os.Exit(1)
	// }

}

// func run(args []string) ([]byte, int, error) {
// 	cmd := exec.Command(args[0], args[1:]...)
// 	stdout, err := cmd.CombinedOutput()
// 	exitCode := cmd.ProcessState.ExitCode()
// 	return stdout, exitCode, err
// }