package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fspath"
	"strings"
)

// TODO: refactor duplicate code
func pathNames(path *fspath.FileSystemPath, workingDir file.File) []string {
	if path.IsAbs() || workingDir == nil {
		if path.AbsolutePath() == "/" {
			return []string{}
		}
		return strings.Split(strings.Trim(path.AbsolutePath(), "/"), "/")
	}

	return strings.Split(strings.Trim(path.Path(), "/"), "/")
}

func pathDirs(path *fspath.FileSystemPath, workingDir file.File) []string {
	if path.IsAbs() || workingDir == nil {
		if path.AbsDir() == "/" {
			return []string{}
		}
		return strings.Split(strings.Trim(path.AbsDir(), "/"), "/")
	}

	return strings.Split(strings.Trim(path.Dir(), "/"), "/")
}
