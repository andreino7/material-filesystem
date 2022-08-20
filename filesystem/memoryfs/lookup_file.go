package memoryfs

import (
	"fmt"
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fspath"
)

func (fs *MemoryFileSystem) lookupDir(pathRoot *inMemoryFile, pathNames []string) (*inMemoryFile, error) {
	return fs.lookupDirWithCreateMissing(pathRoot, pathNames, false)
}

// TODO: refactor
func (fs *MemoryFileSystem) lookupDirWithCreateMissing(pathRoot *inMemoryFile, pathNames []string, createMissing bool) (*inMemoryFile, error) {
	var err error
	for _, currentDir := range pathNames {
		tmp, found := pathRoot.fileMap[currentDir]
		if !found {
			if createMissing {
				tmp, err = fs.createFile(currentDir, true, pathRoot)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, fmt.Errorf("no such file or directory")
			}
		}

		if !tmp.info.IsDirectory() {
			return nil, fmt.Errorf("file is not a directory")
		}

		pathRoot = tmp
	}

	return pathRoot, nil
}

func (fs *MemoryFileSystem) findPathRoot(path *fspath.FileSystemPath, workingDir file.File) (*inMemoryFile, error) {
	if path.IsAbs() || workingDir == nil {
		return fs.root, nil
	}

	return fs.resolveWorkDir(path, workingDir)
}
