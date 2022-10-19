package fixme

import (
	"errors"
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/nabaz-io/nabaz/pkg/fixme/watch"
)

func handleFSEvent(event fsnotify.Event) {
	// do something
	fmt.Printf("event: %v\n", event.String())
}

func Execute(args *Arguements) error {
	fmt.Printf("args: %v\n", args)
	_ = args.Cmdline
	path := args.RepoPath

	fsEvents, err := watch.InitWatch(path)
	if err != nil {
		return err
	}

	for {
		select {
		case event := <-fsEvents:
			handleFSEvent(event)
		}
	}
	return errors.New("HELLO")
}
