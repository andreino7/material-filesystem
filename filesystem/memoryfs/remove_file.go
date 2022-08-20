package memoryfs

import (
	"fmt"
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fspath"
)

func (fs *MemoryFileSystem) RemoveAll(path *fspath.FileSystemPath, workingDir file.File) error {
	return fs.removeFileWithLock(path, workingDir, true)
}

func (fs *MemoryFileSystem) RemoveRegularFile(path *fspath.FileSystemPath, workingDir file.File) error {
	return fs.removeFileWithLock(path, workingDir, false)
}

func (fs *MemoryFileSystem) removeFileWithLock(path *fspath.FileSystemPath, workingDir file.File, isRecursive bool) error {
	// RW lock the fs
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	// find parent directory
	parent, err := fs.lookupParentDirWithCreateMissingDir(path, workingDir, false)
	if err != nil {
		return err
	}

	return fs.removeFile(path.Base(), parent, isRecursive)
}

func (fs *MemoryFileSystem) removeFile(fileName string, parent *inMemoryFile, isRecursive bool) error {
	// check if file exists
	fileToRemove, found := parent.fileMap[fileName]
	if !found {
		return fmt.Errorf("no such file or directory")
	}

	// handle directories
	if fileToRemove.info.IsDirectory() {
		return fs.removeDirectory(fileToRemove, parent, isRecursive)
	}

	// unlink regular file
	fs.unlink(fileToRemove, parent)
	return nil
}

func (fs *MemoryFileSystem) removeDirectory(fileToRemove *inMemoryFile, parent *inMemoryFile, isRecursive bool) error {
	if !isRecursive {
		return fmt.Errorf("file is a directory")
	}

	// remove current file
	fs.unlink(fileToRemove, parent)

	// remove all children
	for _, childFile := range fileToRemove.fileMap {
		if err := fs.removeDirectory(childFile, fileToRemove, true); err != nil {
			return err
		}
	}

	return nil
}

func (fs *MemoryFileSystem) unlink(fileToRemove *inMemoryFile, parent *inMemoryFile) {
	delete(parent.fileMap, fileToRemove.info.Name())
	delete(fileToRemove.fileMap, "..")
	delete(fileToRemove.fileMap, ".")
	fileToRemove.isDeleted = true
}
