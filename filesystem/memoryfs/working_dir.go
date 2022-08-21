package memoryfs

import (
	"fmt"
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fspath"
)

func (fs *MemoryFileSystem) DefaultWorkingDirectory() file.File {
	return fs.root
}

func (fs *MemoryFileSystem) GetDirectory(path *fspath.FileSystemPath, workingDir file.File) (file.File, error) {
	// RLock the fs
	fs.mutex.RLock()
	defer fs.mutex.RUnlock()

	// Find path starting point
	pathRoot, err := fs.findPathRoot(path, workingDir)
	if err != nil {
		return nil, err
	}

	// Find directory
	pathNames := pathNames(path, workingDir)
	return fs.lookupDir(pathRoot, pathNames)
}

// TODO: Test directory deleted
func (fs *MemoryFileSystem) resolveWorkDir(path *fspath.FileSystemPath, workingDir file.File) (*inMemoryFile, error) {
	currentDir, ok := workingDir.(*inMemoryFile)
	if !ok {
		return nil, fmt.Errorf("invalid working directory")
	}

	if currentDir.isDeleted {
		return nil, fmt.Errorf("working directory deleted")
	}

	return currentDir, nil
}
