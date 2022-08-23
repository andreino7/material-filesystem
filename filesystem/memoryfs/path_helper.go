package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"
	"strings"

	"github.com/google/uuid"
)

var invalidFileNames = map[string]bool{
	"..": true,
	".":  true,
	"/":  true,
}

// TODO: refactor duplicate code
func pathNames(path *fspath.FileSystemPath, workingDir file.File) []string {
	if path.Path() == "/" {
		return []string{}
	}

	return strings.Split(strings.Trim(path.Path(), "/"), "/")
}

func pathDirs(path *fspath.FileSystemPath, workingDir file.File) []string {
	if path.Dir() == "/" {
		return []string{}
	}

	return strings.Split(strings.Trim(path.Dir(), "/"), "/")
}

func checkFilePath(path *fspath.FileSystemPath) error {
	return checkFileName(path.Base())
}

func checkFileName(name string) error {
	if _, found := invalidFileNames[name]; found {
		return fserrors.ErrInvalid
	}
	return nil
}

// TODO: use a better way to fix name conflict, for now using UUID is good enough
func generateRandomNameFromBaseName(name string) string {
	return name + "_" + uuid.NewString()
}
