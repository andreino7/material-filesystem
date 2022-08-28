package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"
)

// Remove removes the file located at the specified path.
// The file will be garbage collected once everybody has
// closed it, i.e. has been removed from the open files table.
// This means that is possible to write to a deleted file, but the file will be
// discarded as soon as it is closed.
// This implementation is thread safe.
//
// Returns an error when:
// - The file is a directory
// - The file does not exist
func (fs *MemoryFileSystem) Remove(path *fspath.FileSystemPath) (file.FileInfo, error) {
	return fs.removeFileWithLock(path, false)
}

// RemoveAll removes the file or directory located at the specified path.
// The file will be garbage collected once everybody has
// closed it, i.e. has been removed from the open files table.
// This means that is possible to write to a deleted file, but the file will be
// discarded as soon as it is closed.
// Removing "/" is not supported.
// This implementation is thread safe.
//
// Returns an error when:
// - The file does not exist
func (fs *MemoryFileSystem) RemoveAll(path *fspath.FileSystemPath) (file.FileInfo, error) {
	return fs.removeFileWithLock(path, true)
}

func (fs *MemoryFileSystem) removeFileWithLock(path *fspath.FileSystemPath, isRecursive bool) (file.FileInfo, error) {
	// RW lock the fs
	fs.Lock()
	defer fs.Unlock()

	// find where to remove directory
	pathEnd, err := fs.traverseDirs(path)
	if err != nil {
		return nil, err
	}

	return fs.removeFile(path.Base(), pathEnd, isRecursive)
}

// removeFile removes the file from the fs tree
func (fs *MemoryFileSystem) removeFile(fileName string, pathEnd *inMemoryFile, isRecursive bool) (file.FileInfo, error) {
	// check if file exists
	fileToRemove, found := pathEnd.fileMap[fileName]
	if !found {
		return nil, fserrors.ErrNotExist
	}

	// handle directories
	if fileToRemove.info.fileType == file.Directory {
		return fs.removeDirectory(fileToRemove, pathEnd, isRecursive)
	}

	// unlink regular file (or symlink)
	fs.detachFromParent(fileToRemove)
	// mark file for deletion
	fileToRemove.isDeleted = true
	return fileToRemove.Info(), nil
}

// removeDirectory recursively removes any children file and directories
func (fs *MemoryFileSystem) removeDirectory(fileToRemove *inMemoryFile, parent *inMemoryFile, isRecursive bool) (file.FileInfo, error) {
	if !isRecursive {
		return nil, fserrors.ErrInvalidFileType
	}

	// deleting filesystem root is not supported at the moment
	if fileToRemove == fs.root {
		return nil, fserrors.ErrOperationNotSupported
	}

	// remove current file
	fs.detachFromParent(fileToRemove)
	fileToRemove.isDeleted = true

	// remove all children
	err := fs.visitDir(fileToRemove, func(_ string, file *inMemoryFile) error {
		_, err := fs.removeDirectory(file, fileToRemove, true)
		return err
	})
	if err != nil {
		return nil, err
	}

	return fileToRemove.info, nil
}

func (fs *MemoryFileSystem) detachFromParent(fileToRemove *inMemoryFile) {
	parent := fileToRemove.fileMap[".."]

	delete(parent.fileMap, fileToRemove.info.Name())
	delete(fileToRemove.fileMap, "..")
	delete(fileToRemove.fileMap, ".")
}
