package memoryfs

import (
	"fmt"
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fspath"
	"path/filepath"
)

type onMoveDestFound func(fileToMove *inMemoryFile, dest *inMemoryFile) (file.FileInfo, error)
type onMoveDestNotFound func(fileToMove *inMemoryFile, dest *inMemoryFile, newName string) (file.FileInfo, error)

// Cases:
// srcPath exists
//    - destPath exists
//      - name conflict
//        - error
//      - no conflict
//        - move
//    - destPath not exist
//      - parent exist
//        - name conflict
//			- error
//        - no conflict
//          - move and rename
//      - parent not exist
//        - not recursive
//          - error
//        - recursive
//          - create intermediate dirs
// srcPath not exists
//    - error

// TODO: handle name conflicts as option
// TODO: handle recursive as opttion
func (fs *MemoryFileSystem) Move(srcPath *fspath.FileSystemPath, destPath *fspath.FileSystemPath, workingDir file.File) (file.FileInfo, error) {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	// find the file/directory that needs to be moved
	fileToMove, err := fs.navigateToEndOfPath(srcPath, workingDir, false)
	if err != nil {
		return nil, err
	}

	if fileToMove == fs.root {
		return nil, fmt.Errorf("operation not supported, moving root directory")
	}

	// find last directory in the destination path
	dest, err := fs.navigateToLastDirInPath(destPath, workingDir, true)
	if err != nil {
		return nil, err
	}

	return fs.moveLockFree(fileToMove, dest, destPath.Base())
}

// This method should be called only if the caller has already acquired a lock
func (fs *MemoryFileSystem) moveLockFree(fileToMove *inMemoryFile, dest *inMemoryFile, finalDestName string) (file.FileInfo, error) {
	if fileToMove.info.isDirectory {
		return fs.moveDirectory(fileToMove, dest, finalDestName)
	}
	return fs.moveRegularFile(fileToMove, dest, finalDestName)
}

func (fs *MemoryFileSystem) moveDirectory(fileToMove *inMemoryFile, dest *inMemoryFile, finalDestName string) (file.FileInfo, error) {
	return fs.doMove(fileToMove, dest, finalDestName, fs.moveDirectoryToExistingDestination, fs.renameAndMoveDirectory)
}

func (fs *MemoryFileSystem) renameAndMoveDirectory(fileToMove *inMemoryFile, dest *inMemoryFile, newName string) (file.FileInfo, error) {
	return fs.renameAndMoveRegularFile(fileToMove, dest, newName)
}

// This method moves the source directory to an existing destination.
// If destination is a directory, the source directory and destination directory are merged.
// if destination is a regular file, the source directory is move to the destination's parent.
func (fs *MemoryFileSystem) moveDirectoryToExistingDestination(fileToMove *inMemoryFile, dest *inMemoryFile) (file.FileInfo, error) {
	if dest.info.isDirectory {
		return fs.mergeDirectories(fileToMove, dest)
	}

	// move to dest parent dir, and rename
	return fs.renameAndMoveRegularFile(fileToMove, dest.fileMap[".."], fileToMove.info.Name())
}

// This method merges two directories and recursively all the subdirectories.
// If in the destination directory there is no directory with same name as the source directory,
// the source directory is simply moved to the new location.
// If in the destination directory there is a directory with the same name as the source directory,
// all the files in the source directory are moved the destination directory and the source directory
// is removed.
func (fs *MemoryFileSystem) mergeDirectories(dirToMove *inMemoryFile, dest *inMemoryFile) (file.FileInfo, error) {
	finalDest, found := dest.fileMap[dirToMove.info.Name()]
	if !found {
		return fs.renameAndMoveDirectory(dirToMove, dest, dirToMove.info.Name())
	}

	// This is the more complex case: recursively move every file to destination directory
	for fileName, fileToMove := range dirToMove.fileMap {
		if fileName == "." || fileName == ".." {
			continue
		}

		if shouldMergeSubDirectories(fileToMove, finalDest) {
			if _, err := fs.mergeDirectories(fileToMove, finalDest); err != nil {
				return nil, err
			}
		} else {
			if _, err := fs.moveLockFree(fileToMove, finalDest, fileName); err != nil {
				return nil, err
			}
		}
	}

	// TODO: document that working dir will be reset when this happens
	fs.removeDirectory(dirToMove, dirToMove.fileMap[".."], true)
	return finalDest.info, nil
}

func shouldMergeSubDirectories(fileToMove *inMemoryFile, dest *inMemoryFile) bool {
	if !fileToMove.info.isDirectory {
		return false
	}

	_, found := dest.fileMap[fileToMove.info.Name()]
	return found
}

func (fs *MemoryFileSystem) moveRegularFile(fileToMove *inMemoryFile, dest *inMemoryFile, finalDestName string) (file.FileInfo, error) {
	return fs.doMove(fileToMove, dest, finalDestName, fs.moveRegularFileToExistingDestination, fs.renameAndMoveRegularFile)
}

// This method moves a file to an existing destination.
// If the destination is a directory, the file is added to the directory.
// If the destination is a file, the source file is moved to the destination's parent
func (fs *MemoryFileSystem) moveRegularFileToExistingDestination(fileToMove *inMemoryFile, dest *inMemoryFile) (file.FileInfo, error) {
	finalDir := dest
	newName := fileToMove.info.Name()

	// Moving to a file, i.e. need to rename
	if !dest.info.isDirectory {
		finalDir = dest.fileMap[".."]
		newName = dest.Info().Name()
	}

	return fs.renameAndMoveRegularFile(fileToMove, finalDir, newName)
}

// This method unlinks the source file from the original location and links it to the new location and renames it
// to the given name.
// If there's a name the conflict the source file is automatically renamed.
func (fs *MemoryFileSystem) renameAndMoveRegularFile(fileToMove *inMemoryFile, dest *inMemoryFile, newName string) (file.FileInfo, error) {
	// detach from parent dir
	fs.unlink(fileToMove)

	// check for name conflicts
	finalName := newName
	if _, found := dest.fileMap[finalName]; found {
		finalName = generateRandomNameFromBaseName(finalName)
	}

	// rename
	fileToMove.info.setAbsolutePath(filepath.Join(dest.info.AbsolutePath(), finalName))

	// attach to new dir
	fs.linkToParent(fileToMove, dest)

	return fileToMove.info, nil
}

func (fs *MemoryFileSystem) doMove(fileToMove *inMemoryFile, dest *inMemoryFile, finalDestName string, onFound onMoveDestFound, onNotFound onMoveDestNotFound) (file.FileInfo, error) {
	// check if dest file exists already
	finalDest, found := dest.fileMap[finalDestName]
	if found {
		return onFound(fileToMove, finalDest)
	} else {
		// rename file
		return onNotFound(fileToMove, dest, finalDestName)
	}
}
