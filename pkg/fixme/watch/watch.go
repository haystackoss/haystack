package watch

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func isWatchedFile(file string) bool {
	validExt := []string{
		".tmpl",
		".tpl",
		".go",
		".py",
		".c",
		".cpp",
		".h",
		".hpp",
		".sh",
	}
	for _, ext := range validExt {
		if strings.HasSuffix(file, ext) {
			return true
		}
	}

	return false
}

func watchFolder(path string, fsEvents chan<- fsnotify.Event) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	err = watcher.Add(path) // Add more stuff
	if err != nil {
		panic(err)
	}

	folderExists := make(chan bool, 1)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				fsEvents <- event
				if event.Op == fsnotify.Remove && event.Name == path {
					folderExists <- false
					return
				}
			case err := <-watcher.Errors:
				println("error:", err)
			}
		}
	}()

	<-folderExists
	return
}

func isIgnoredFolder(path string) bool {
	// TODO: i.e .git, node_modules, pycache, etc
	ignoredFolders := []string{
		".git",
		"node_modules",
		"__pycache__",
		".idea",
		".vscode",
		".cache",
		".pytest_cache",
		".mypy_cache",
		".tox",
		".eggs",
		".venv",
		".env",
	}
	basePath := filepath.Base(path)
	for _, folder := range ignoredFolders {
		if basePath == folder {
			return true
		}
	}
	return false
}

// InitWatch inits a watch on all recursive folders under the path, its called a watch initiation -
// because if a folder is inevitably created, then it won't be covered here.
func InitWatch(path string) (chan fsnotify.Event, error) {
	root := path
	fsEvents := make(chan fsnotify.Event, 2)
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if len(path) > 1 && strings.HasPrefix(filepath.Base(path), ".") {
				return filepath.SkipDir
			}

			if isIgnoredFolder(path) {
				return filepath.SkipDir
			}

			watchFolder(path, fsEvents)
		}
		return err
	})
	return fsEvents, nil
}
