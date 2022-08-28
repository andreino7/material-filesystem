package memoryfs

import (
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"
	"strings"

	"github.com/google/uuid"
)

// Restricted file names
var invalidFileNames = map[string]bool{
	"..": true,
	".":  true,
	"/":  true,
}

// pathDirs returns a list of "dirs" contained in the path
func pathDirs(path *fspath.FileSystemPath) []string {
	if path.Dir() == "/" {
		return []string{}
	}

	return strings.Split(strings.Trim(path.Dir(), "/"), "/")
}

// checkFilePath checks if a path is valid
func checkFilePath(path *fspath.FileSystemPath) error {
	return checkFileName(path.Base())
}

// checkFileName checks if a name is valid
func checkFileName(name string) error {
	if _, found := invalidFileNames[name]; found {
		return fserrors.ErrInvalid
	}
	return nil
}

// generateRandomNameFromBaseName generates unique file name using UUID
//
// TODO: use a better way to fix name conflict, for now using UUID is good enough
func generateRandomNameFromBaseName(name string) string {
	return name + "_" + uuid.NewString()
}
