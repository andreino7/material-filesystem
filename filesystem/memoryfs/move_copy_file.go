package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"
	"path/filepath"
)

type onMoveOrCopyDestFound func(fileToMove *inMemoryFile, dest *inMemoryFile, isCopy bool) (*inMemoryFile, error)
type onMoveOrCopyDestNotFound func(fileToMove *inMemoryFile, dest *inMemoryFile, newName string, isCopy bool) (*inMemoryFile, error)

// Move moves (renames) srcPath to destPath and creates
// any parent directories.
// If destPath exists and is not a directory, the
// "moved" file is automatically renamed to a unique name.
// If destPath exists and is a directory, the
// directories are merged and any name conflict is fixed.
// Move stops at the first error encountered.
// Moving "/" is not supported.
// This implementation is thread safe
//
// Returns an error when:
// - srcPath does not exist
// - the new file name is invalid
//
// TODO: handle name conflicts as option
// TODO: handle create parent dirs as opttion
func (fs *MemoryFileSystem) Move(srcPath *fspath.FileSystemPath, destPath *fspath.FileSystemPath) (file.FileInfo, error) {
	return fs.moveOrCopy(srcPath, destPath, false)
}

// Copy copies srcPath to destPath and creates
// any parent directories.
// If destPath exists and is not a directory, the
// "copied" file is automatically renamed to a unique name.
// If destPath exists and is a directory, the
// directories are merged and any name conflict is fixed.
// Copy stops at the first error encountered.
// Limitation: Copying "/" is not supported.
// This implementation is thread safe
//
// Returns an error when:
// - srcPath does not exist
// - the new file name is invalid
//
// TODO: handle name conflicts as option
// TODO: handle create parent dirs as opttion
func (fs *MemoryFileSystem) Copy(srcPath *fspath.FileSystemPath, destPath *fspath.FileSystemPath) (file.FileInfo, error) {
	return fs.moveOrCopy(srcPath, destPath, true)
}

func (fs *MemoryFileSystem) moveOrCopy(srcPath *fspath.FileSystemPath, destPath *fspath.FileSystemPath, isCopy bool) (file.FileInfo, error) {
	fs.Lock()
	defer fs.Unlock()

	// find the file/directory that needs to be moved/copied
	fileToMove, err := fs.traverseToBaseWithSkipLastLink(srcPath, !isCopy)
	if err != nil {
		return nil, err
	}

	if fileToMove == fs.root {
		return nil, fserrors.ErrOperationNotSupported
	}

	// find last directory in the destination path
	dest, err := fs.traverseDirsAndCreateParentDirs(destPath)
	if err != nil {
		return nil, err
	}

	newFile, err := fs.moveOrCopyLockFree(fileToMove, dest, destPath.Base(), isCopy)
	if err != nil {
		return nil, err
	}
	return newFile.info, nil
}

// This method should be called only if the caller has already acquired a lock
func (fs *MemoryFileSystem) moveOrCopyLockFree(fileToMove *inMemoryFile, dest *inMemoryFile, finalDestName string, isCopy bool) (*inMemoryFile, error) {
	if fileToMove.info.fileType == file.Directory {
		return fs.moveOrCopyDirectory(fileToMove, dest, finalDestName, isCopy)
	}
	return fs.moveOrCopyRegularFile(fileToMove, dest, finalDestName, isCopy)
}

func (fs *MemoryFileSystem) moveOrCopyDirectory(fileToMove *inMemoryFile, dest *inMemoryFile, finalDestName string, isCopy bool) (*inMemoryFile, error) {
	return fs.doMoveOrCopy(fileToMove, dest, finalDestName, fs.moveOrCopyDirectoryToExistingDestination, fs.renameAndMoveOrCopyDirectory, isCopy)
}

func (fs *MemoryFileSystem) renameAndMoveOrCopyDirectory(fileToMove *inMemoryFile, dest *inMemoryFile, newName string, isCopy bool) (*inMemoryFile, error) {
	return fs.renameAndMoveOrCopy(fileToMove, dest, newName, isCopy)
}

// moveOrCopyDirectoryToExistingDestination moves/copies the source directory to an existing destination.
// If destination is a directory, the source directory and destination directory are merged.
// if destination is a regular file, the source directory is moved/copied to the destination's parent and renamed.
func (fs *MemoryFileSystem) moveOrCopyDirectoryToExistingDestination(fileToMove *inMemoryFile, dest *inMemoryFile, isCopy bool) (*inMemoryFile, error) {
	if dest.info.fileType == file.Directory {
		return fs.mergeDirectories(fileToMove, dest, isCopy)
	}

	// move/copy to dest parent dir, and rename
	return fs.renameAndMoveOrCopy(fileToMove, dest.fileMap[".."], fileToMove.info.Name(), isCopy)
}

// mergeDirectories merges two directories and recursively all the subdirectories.
// If in the destination directory there is no directory with same name as the source directory,
// the source directory is simply moved/copied to the new location.
// If in the destination directory there is a directory with the same name as the source directory,
// all the files in the source directory are moved/copied to the destination directory and, in case of a move,
// the source directory is removed.
func (fs *MemoryFileSystem) mergeDirectories(dirToMove *inMemoryFile, dest *inMemoryFile, isCopy bool) (*inMemoryFile, error) {
	finalDest, found := dest.fileMap[dirToMove.info.Name()]
	if !found {
		return fs.renameAndMoveOrCopyDirectory(dirToMove, dest, dirToMove.info.Name(), isCopy)
	}

	// This is the more complex case: recursively move/copy every file to destination directory
	err := fs.visitDir(dirToMove, func(fileName string, fileToMove *inMemoryFile) error {
		var err error
		if shouldMergeSubDirectories(fileToMove, finalDest) {
			_, err = fs.mergeDirectories(fileToMove, finalDest, isCopy)
		} else {
			_, err = fs.moveOrCopyLockFree(fileToMove, finalDest, fileName, isCopy)
		}
		return err
	})

	if err != nil {
		return nil, err
	}

	if !isCopy {
		fs.removeDirectory(dirToMove, true)
	}
	return finalDest, nil
}

// shouldMergeSubDirectories returns true if source is a directory and destination
// contains a directory with the same name.
func shouldMergeSubDirectories(fileToMove *inMemoryFile, dest *inMemoryFile) bool {
	if fileToMove.info.fileType != file.Directory {
		return false
	}

	_, found := dest.fileMap[fileToMove.info.Name()]
	return found
}

func (fs *MemoryFileSystem) moveOrCopyRegularFile(fileToMove *inMemoryFile, dest *inMemoryFile, finalDestName string, isCopy bool) (*inMemoryFile, error) {
	return fs.doMoveOrCopy(fileToMove, dest, finalDestName, fs.moveOrCopyRegularFileToExistingDestination, fs.renameAndMoveOrCopy, isCopy)
}

// moveOrCopyRegularFileToExistingDestination moves/copies a file to an existing destination.
// If the destination is a directory, the file is moved/copied to the directory.
// If the destination is a file, the source file is move/copied to the destination's parent and renamed
func (fs *MemoryFileSystem) moveOrCopyRegularFileToExistingDestination(fileToMove *inMemoryFile, dest *inMemoryFile, isCopy bool) (*inMemoryFile, error) {
	finalDir := dest
	newName := fileToMove.info.Name()

	// Moving to a file, i.e. need to rename
	if dest.info.fileType != file.Directory {
		finalDir = dest.fileMap[".."]
		newName = dest.Info().Name()
	}

	return fs.renameAndMoveOrCopy(fileToMove, finalDir, newName, isCopy)
}

// renameAndMoveOrCopy In case of "Move", removes the source file from the original location
// and attaches it to the new location and renames it to the given name.
// In case of "Copy", copies the source file from the original location
// and attaches it to the new location and renames it to the given name.
// If there's a name the conflict the source file is automatically renamed.
func (fs *MemoryFileSystem) renameAndMoveOrCopy(fileToMove *inMemoryFile, dest *inMemoryFile, newName string, isCopy bool) (*inMemoryFile, error) {
	// check for name conflicts
	finalName := newName
	if _, found := dest.fileMap[finalName]; found {
		finalName = generateRandomNameFromBaseName(finalName)
	}

	// check if new name is valid
	if err := checkFileName(finalName); err != nil {
		return nil, err
	}

	newAbsPath := filepath.Join(dest.info.AbsolutePath(), finalName)

	var result *inMemoryFile
	var err error
	if isCopy {
		result, err = fs.copyFile(fileToMove, newAbsPath)
	} else {
		result, err = fs.moveFile(fileToMove, newAbsPath)
	}
	if err != nil {
		return nil, err
	}

	// attach to new dir
	fs.attachToParent(result, dest)

	return result, nil
}

// moveFile detaches the file from the original parent and uptades the absolute path
func (fs *MemoryFileSystem) moveFile(fileToMove *inMemoryFile, newAbsPath string) (*inMemoryFile, error) {

	// detach from parent dir
	fs.detachFromParent(fileToMove)
	// Update absolute path
	return fs.updatePaths(fileToMove, newAbsPath)
}

// updatePaths - updates the path of the given file. If it's a directory
// recursively update every children
func (fs *MemoryFileSystem) updatePaths(fileToUpdate *inMemoryFile, newAbsPath string) (*inMemoryFile, error) {
	fileToUpdate.info.absolutePath = newAbsPath

	if fileToUpdate.info.fileType == file.Directory {
		err := fs.visitDir(fileToUpdate, func(fileName string, child *inMemoryFile) error {
			fs.updatePaths(child, filepath.Join(newAbsPath, fileName))
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return fileToUpdate, nil
}

// copyFile creates a copy of the original file.
// If the file is a directory recursively copies every file in it
func (fs *MemoryFileSystem) copyFile(fileToMove *inMemoryFile, newAbsPath string) (*inMemoryFile, error) {
	newFile := newInMemoryFile(newAbsPath, fileToMove.info.fileType)

	if fileToMove.info.fileType == file.Directory {
		err := fs.visitDir(fileToMove, func(fileName string, child *inMemoryFile) error {
			fs.renameAndMoveOrCopy(child, newFile, fileName, true)
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else if fileToMove.info.fileType == file.RegularFile {
		newDataArr := make([]byte, fileToMove.data.Size())
		copy(newDataArr, fileToMove.data.data)
		newData := &inMemoryFileData{
			data: newDataArr,
		}
		newFile.data = newData
	} else {
		newFile.link = fileToMove.link
	}

	return newFile, nil
}

func (fs *MemoryFileSystem) doMoveOrCopy(fileToMove *inMemoryFile, dest *inMemoryFile, finalDestName string, onFound onMoveOrCopyDestFound, onNotFound onMoveOrCopyDestNotFound, isCopy bool) (*inMemoryFile, error) {
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
