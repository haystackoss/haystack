package watcher

var validExtentions = []string{
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

var resourceFilesExt = []string{
	".json",
	".xml",
	".yml",
	".yaml",
	".conf",
	".config",
	".toml",
}

var ignoredFolders = []string{
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
	".nabazgit",
}
