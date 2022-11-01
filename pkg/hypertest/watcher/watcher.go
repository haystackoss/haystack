package watcher

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/nabaz-io/nabaz/pkg/hypertest/limit"
)

type Watcher struct {
	rootPath         string
	FileSystemEvents chan fsnotify.Event
	Errors           chan error
}

func isWatchedFile(file string) bool {
	for _, ext := range validExtentions {
		if strings.HasSuffix(file, ext) {
			return true
		}
	}
	for _, ext := range resourceFilesExt {
		if strings.HasSuffix(file, ext) {
			return true
		}
	}

	return false
}

func (w *Watcher) WatchFolder(path string) {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		limit.InitLimit()
		watcher, err = fsnotify.NewWatcher()
		if err != nil {
			panic(err)
		}
	}
	defer watcher.Close()

	folderExists := make(chan bool, 1)
	go func() {
		defer close(folderExists)
		for {
			select {
			case event := <-watcher.Events:
				if !isWatchedFile(event.Name) {
					continue
				}

				w.FileSystemEvents <- event
				if event.Op == fsnotify.Remove && event.Name == path {
					return
				}
			case err := <-watcher.Errors:
				w.Errors <- err
			}
		}
	}()

	err = watcher.Add(path)
	if err != nil {
		panic(err)
	}
	<-folderExists
}

func isIgnoredFolder(folderName string) bool {
	// TODO: i.e .git, node_modules, pycache, etc
	for _, folder := range ignoredFolders {
		if folderName == folder {
			return true
		}
	}
	return false
}

// InitWatch inits a watch on all recursive folders under the path, its called a watch initiation -
// because if a folder is inevitably created, then it won't be covered here.
func (w *Watcher) initWatch(root string) {
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if len(path) > 1 && strings.HasPrefix(filepath.Base(path), ".") {
				return filepath.SkipDir
			}

			if isIgnoredFolder(info.Name()) {
				return filepath.SkipDir
			}

			go w.WatchFolder(path)
		}
		return err
	})
}

func NewWatcher(rootPath string) *Watcher {
	w := &Watcher{
		rootPath:         rootPath,
		FileSystemEvents: make(chan fsnotify.Event, 2),
		Errors:           make(chan error, 100),
	}

	w.initWatch(w.rootPath)

	return w
}
