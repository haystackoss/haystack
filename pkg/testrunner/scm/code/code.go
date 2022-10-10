package code

import (
	"io/ioutil"
)

// CodeDirectory is the directory where the code is stored.
type CodeDirectory struct {
	cache map[string][]byte
	// Directory is the directory where the code is stored.
	Path string
}

// NewCodeDirectory creates a new CodeDirectory.
func NewCodeDirectory(path string) *CodeDirectory {
	cache := make(map[string][]byte)

	return &CodeDirectory{
		cache: cache,
		Path:  path,
	}
}

func (c *CodeDirectory) GetFileContent(path string) ([]byte, error) {
	if val, ok := c.cache[path]; ok {
		return val, nil
	}


	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	c.cache[path] = fileContent

	return fileContent, nil
}
