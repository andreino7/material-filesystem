package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"
	"path/filepath"
)

type onMoveDestFound func(fileToMove *inMemoryFile, dest *inMemoryFile, isCopy bool) (*inMemoryFile, error)
type onMoveDestNotFound func(fileToMove *inMemoryFile, dest *inMemoryFile, newName string, isCopy bool) (*inMemoryFile, error)

// TODO: handle name conflicts as option
// TODO: handle recursive as opttion
func (fs *MemoryFileSystem) Move(srcPath *fspath.FileSystemPath, destPath *fspath.FileSystemPath, workingDir file.File) (file.FileInfo, error) {
	return fs.moveOrCopy(srcPath, destPath, workingDir, false)
}

func (fs *MemoryFileSystem) Copy(srcPath *fspath.FileSystemPath, destPath *fspath.FileSystemPath, workingDir file.File) (file.FileInfo, error) {
	return fs.moveOrCopy(srcPath, destPath, workingDir, true)
}

func (fs *MemoryFileSystem) moveOrCopy(srcPath *fspath.FileSystemPath, destPath *fspath.FileSystemPath, workingDir file.File, isCopy bool) (file.FileInfo, error) {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	// find the file/directory that needs to be moved
	fileToMove, err := fs.navigateToEndOfPath(srcPath, workingDir, false)
	if err != nil {
		return nil, err
	}

	if fileToMove == fs.root && !isCopy {
		return nil, fserrors.ErrOperationNotSupported
	}

	// find last directory in the destination path
	dest, err := fs.navigateToLastDirInPath(destPath, workingDir, true)
	if err != nil {
		return nil, err
	}

	newFile, err := fs.moveLockFree(fileToMove, dest, destPath.Base(), isCopy)
	if err != nil {
		return nil, err
	}
	return newFile.info, nil
}

// This method should be called only if the caller has already acquired a lock
func (fs *MemoryFileSystem) moveLockFree(fileToMove *inMemoryFile, dest *inMemoryFile, finalDestName string, isCopy bool) (*inMemoryFile, error) {
	if fileToMove.info.isDirectory {
		return fs.moveDirectory(fileToMove, dest, finalDestName, isCopy)
	}
	return fs.moveRegularFile(fileToMove, dest, finalDestName, isCopy)
}

func (fs *MemoryFileSystem) moveDirectory(fileToMove *inMemoryFile, dest *inMemoryFile, finalDestName string, isCopy bool) (*inMemoryFile, error) {
	return fs.doMove(fileToMove, dest, finalDestName, fs.moveDirectoryToExistingDestination, fs.renameAndMoveDirectory, isCopy)
}

func (fs *MemoryFileSystem) renameAndMoveDirectory(fileToMove *inMemoryFile, dest *inMemoryFile, newName string, isCopy bool) (*inMemoryFile, error) {
	return fs.renameAndMoveRegularFile(fileToMove, dest, newName, isCopy)
}

// This method moves the source directory to an existing destination.
// If destination is a directory, the source directory and destination directory are merged.
// if destination is a regular file, the source directory is move to the destination's parent.
func (fs *MemoryFileSystem) moveDirectoryToExistingDestination(fileToMove *inMemoryFile, dest *inMemoryFile, isCopy bool) (*inMemoryFile, error) {
	if dest.info.isDirectory {
		return fs.mergeDirectories(fileToMove, dest, isCopy)
	}

	// move to dest parent dir, and rename
	return fs.renameAndMoveRegularFile(fileToMove, dest.fileMap[".."], fileToMove.info.Name(), isCopy)
}

// This method merges two directories and recursively all the subdirectories.
// If in the destination directory there is no directory with same name as the source directory,
// the source directory is simply moved to the new location.
// If in the destination directory there is a directory with the same name as the source directory,
// all the files in the source directory are moved the destination directory and the source directory
// is removed.
func (fs *MemoryFileSystem) mergeDirectories(dirToMove *inMemoryFile, dest *inMemoryFile, isCopy bool) (*inMemoryFile, error) {
	finalDest, found := dest.fileMap[dirToMove.info.Name()]
	if !found {
		return fs.renameAndMoveDirectory(dirToMove, dest, dirToMove.info.Name(), isCopy)
	}

	// This is the more complex case: recursively move every file to destination directory
	err := fs.walk(dirToMove, func(fileName string, fileToMove *inMemoryFile) error {
		if shouldMergeSubDirectories(fileToMove, finalDest) {
			if _, err := fs.mergeDirectories(fileToMove, finalDest, isCopy); err != nil {
				return err
			}
		} else {
			if _, err := fs.moveLockFree(fileToMove, finalDest, fileName, isCopy); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	if !isCopy {
		// TODO: document that working dir will be reset when this happens
		fs.removeDirectory(dirToMove, dirToMove.fileMap[".."], true)
	}
	return finalDest, nil
}

func shouldMergeSubDirectories(fileToMove *inMemoryFile, dest *inMemoryFile) bool {
	if !fileToMove.info.isDirectory {
		return false
	}

	_, found := dest.fileMap[fileToMove.info.Name()]
	return found
}

func (fs *MemoryFileSystem) moveRegularFile(fileToMove *inMemoryFile, dest *inMemoryFile, finalDestName string, isCopy bool) (*inMemoryFile, error) {
	return fs.doMove(fileToMove, dest, finalDestName, fs.moveRegularFileToExistingDestination, fs.renameAndMoveRegularFile, isCopy)
}

// This method moves a file to an existing destination.
// If the destination is a directory, the file is added to the directory.
// If the destination is a file, the source file is moved to the destination's parent
func (fs *MemoryFileSystem) moveRegularFileToExistingDestination(fileToMove *inMemoryFile, dest *inMemoryFile, isCopy bool) (*inMemoryFile, error) {
	finalDir := dest
	newName := fileToMove.info.Name()

	// Moving to a file, i.e. need to rename
	if !dest.info.isDirectory {
		finalDir = dest.fileMap[".."]
		newName = dest.Info().Name()
	}

	return fs.renameAndMoveRegularFile(fileToMove, finalDir, newName, isCopy)
}

// This method unlinks the source file from the original location and links it to the new location and renames it
// to the given name.
// If there's a name the conflict the source file is automatically renamed.
func (fs *MemoryFileSystem) renameAndMoveRegularFile(fileToMove *inMemoryFile, dest *inMemoryFile, newName string, isCopy bool) (*inMemoryFile, error) {
	// check for name conflicts
	finalName := newName
	if _, found := dest.fileMap[finalName]; found {
		finalName = generateRandomNameFromBaseName(finalName)
	}

	newAbsPath := filepath.Join(dest.info.AbsolutePath(), finalName)

	var result *inMemoryFile
	var err error
	if isCopy {
		result, err = fs.copyFile(fileToMove, newAbsPath)
	} else {
		result = fs.moveFile(fileToMove, newAbsPath)
	}
	if err != nil {
		return nil, err
	}

	// attach to new dir
	fs.attachToParent(result, dest)

	return result, nil
}

func (fs *MemoryFileSystem) moveFile(fileToMove *inMemoryFile, newAbsPath string) *inMemoryFile {
	// detach from parent dir
	fs.detachFromParent(fileToMove)
	fileToMove.info.absolutePath = newAbsPath
	return fileToMove
}

func (fs *MemoryFileSystem) copyFile(fileToMove *inMemoryFile, newAbsPath string) (*inMemoryFile, error) {
	newFile := newInMemoryFile(newAbsPath, fileToMove.info.IsDirectory())

	if fileToMove.info.isDirectory {
		err := fs.walk(fileToMove, func(fileName string, child *inMemoryFile) error {
			fs.renameAndMoveRegularFile(child, newFile, fileName, true)
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		var newDataArr []byte
		copy(newDataArr, fileToMove.data.data)
		newData := &inMemoryFileData{
			data: newDataArr,
		}
		newFile.data = newData
	}

	return newFile, nil
}

func (fs *MemoryFileSystem) doMove(fileToMove *inMemoryFile, dest *inMemoryFile, finalDestName string, onFound onMoveDestFound, onNotFound onMoveDestNotFound, isCopy bool) (*inMemoryFile, error) {
	// check if dest file exists already
	finalDest, found := dest.fileMap[finalDestName]
	if found {
		// check if same filex
		if finalDest == fileToMove {
			return nil, fserrors.ErrSameFile
		}
		return onFound(fileToMove, finalDest, isCopy)
	} else {
		// rename file
		return onNotFound(fileToMove, dest, finalDestName, isCopy)
	}
}
