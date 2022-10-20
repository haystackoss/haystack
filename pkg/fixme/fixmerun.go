package fixme

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/nabaz-io/nabaz/pkg/fixme/watcher"
)

// handleFSCreate assumes that the file was just created and not already watched.
func handleFSCreate(w *watcher.Watcher, event fsnotify.Event) {

	info, err := os.Lstat(event.Name)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	if info.IsDir() {
		w.WatchFolder(event.Name)
	}
}

func handleFSEvent(w *watcher.Watcher, event fsnotify.Event) {
	// do something
	switch event.Op {
	case fsnotify.Create:
		handleFSCreate(w, event)
	}
}

func Execute(args *Arguements) error {
	cmdline := args.Cmdline
	path := args.RepoPath
	_ = cmdline

	w := watcher.NewWatcher(path)

	for {
		select {
		case event := <-w.FileSystemEvents:
			handleFSEvent(w, event)
		case err := <-w.Errors:
			fmt.Printf("error: %v\n", err)
		}
	}

	return nil
}
