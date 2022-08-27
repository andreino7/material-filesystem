package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"
	"material/filesystem/filesystem/user"
)

// TODO: handle deleting working dir
func (fs *MemoryFileSystem) Remove(path *fspath.FileSystemPath, user user.User) (file.FileInfo, error) {
	return fs.removeFileWithLock(path, false, user)
}

func (fs *MemoryFileSystem) RemoveAll(path *fspath.FileSystemPath, user user.User) (file.FileInfo, error) {
	return fs.removeFileWithLock(path, true, user)
}

func (fs *MemoryFileSystem) removeFileWithLock(path *fspath.FileSystemPath, isRecursive bool, user user.User) (file.FileInfo, error) {
	// RW lock the fs
	fs.Lock()
	defer fs.Unlock()

	// find where to remove directory
	pathEnd, err := fs.traverseToDir(path, user)
	if err != nil {
		return nil, err
	}

	return fs.removeFile(path.Base(), pathEnd, isRecursive)
}

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
	fileToRemove.isDeleted = true
	return fileToRemove.Info(), nil
}

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
	err := fs.walk(fileToRemove, func(_ string, file *inMemoryFile) error {
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
