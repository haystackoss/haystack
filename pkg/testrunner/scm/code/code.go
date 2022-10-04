package code

import (
	"io/ioutil"
	"path/filepath"
)

/*
class LocalCodeDirectory:

    def __init__(self, code_path: str = "."):
        self._logger = logging.getLogger(self.__class__.__name__)
        self._cache = dict()
        self.local_code_directory_path = code_path
        self._validate()

    def _validate(self):
        if not os.path.isdir(self.local_code_directory_path):
            raise click.UsageError(f"Code path \"{self.local_code_directory_path}\" is not a directory."
                                   f"Please pass a valid path to the code directory.")
*/

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

func (c *CodeDirectory) GetFileContent(filePath string) ([]byte, error) {
	if val, ok := c.cache[filePath]; ok {
		return val, nil
	}

	filePath = filepath.Join(c.Path, filePath)

	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	c.cache[filePath] = fileContent

	return fileContent, nil
}
