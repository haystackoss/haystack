package paths

import "os"


func TempDir() string {
	tmpdir := os.TempDir()
	if tmpdir == "" {
		nomedir, err := os.UserHomeDir()
		if err != nil {
			tmpdir = "."
		} else {
			tmpdir = nomedir
		}
	}
	return tmpdir
}

func JunitXMLName() string {
	return "nabaz-junit.xml"
}

func JunitXMLPath() string {
	return TempDir() + "/nabaz-junit.xml"
}