package code

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/nabaz-io/nabaz/pkg/testrunner/scm/git/local"
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
-
    def _is_dotgit_provided(self) -> bool:
        return os.path.isdir(os.path.join(self.local_code_directory_path, ".git"))

    def find_local_git_repo(self) -> Union[LocalGitRepo, None]:
        if not self._is_dotgit_provided():
            return None

        return LocalGitRepo(repo_path=self.local_code_directory_path)
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

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}

func (c *CodeDirectory) isDotGitProvided() bool {
	return isDirectory(filepath.Join(c.Path, ".git"))
}

func (c *CodeDirectory) FindLocalGitRepo() (*local.LocalGitRepo, error) {
	if !c.isDotGitProvided() {
		return nil, nil
	}

	return local.NewLocalGitRepo(c.Path)
}
