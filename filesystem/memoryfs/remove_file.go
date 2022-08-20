package memoryfs

import (
	"fmt"
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fspath"
)

// TODO: handle deleting working dir
func (fs *MemoryFileSystem) RemoveDirectory(path *fspath.FileSystemPath, workingDir file.File) (file.FileInfo, error) {
	return fs.removeFileWithLock(path, workingDir, true)
}

func (fs *MemoryFileSystem) RemoveRegularFile(path *fspath.FileSystemPath, workingDir file.File) (file.FileInfo, error) {
	return fs.removeFileWithLock(path, workingDir, false)
}

func (fs *MemoryFileSystem) removeFileWithLock(path *fspath.FileSystemPath, workingDir file.File, isRecursive bool) (file.FileInfo, error) {
	// RW lock the fs
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	// find where to remove directory
	pathEnd, err := fs.lookupPathEndWithCreateMissingDir(path, workingDir, false)
	if err != nil {
		return nil, err
	}

	return fs.removeFile(path.Base(), pathEnd, isRecursive)
}

func (fs *MemoryFileSystem) removeFile(fileName string, pathEnd *inMemoryFile, isRecursive bool) (file.FileInfo, error) {
	// check if file exists
	fileToRemove, found := pathEnd.fileMap[fileName]
	if !found {
		return nil, fmt.Errorf("no such file or directory")
	}

	// handle directories
	if fileToRemove.info.IsDirectory() {
		return fs.removeDirectory(fileToRemove, pathEnd, isRecursive)
	}

	// unlink regular file
	fs.unlink(fileToRemove)
	return fileToRemove.Info(), nil
}

func (fs *MemoryFileSystem) removeDirectory(fileToRemove *inMemoryFile, parent *inMemoryFile, isRecursive bool) (file.FileInfo, error) {
	if !isRecursive {
		return nil, fmt.Errorf("file is a directory")
	}

	// deleting filesystem root is not supported at the moment
	if fileToRemove == fs.root {
		return nil, fmt.Errorf("operation not allowed: deleting filesystem root")
	}

	// remove current file
	fs.unlink(fileToRemove)

	// remove all children
	for _, nextFile := range fileToRemove.fileMap {
		if _, err := fs.removeDirectory(nextFile, fileToRemove, true); err != nil {
			return nil, err
		}
	}

	return fileToRemove.info, nil
}

func (fs *MemoryFileSystem) unlink(fileToRemove *inMemoryFile) {
	parent := fileToRemove.fileMap[".."]

	delete(parent.fileMap, fileToRemove.info.Name())
	delete(fileToRemove.fileMap, "..")
	delete(fileToRemove.fileMap, ".")
	fileToRemove.isDeleted = true
}
