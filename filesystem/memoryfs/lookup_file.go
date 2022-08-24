package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"
)

type onFileNotFoundFn func(string, *inMemoryFile) (*inMemoryFile, error)
type onFileFoundFn func(*inMemoryFile) (*inMemoryFile, error)

const MAX_LINK_DEPTH = 40

var (
	// return error if file not found
	errorOnNotFoundFn = func(filename string, parent *inMemoryFile) (*inMemoryFile, error) {
		return nil, fserrors.ErrNotExist
	}
	// do nothing
	noopOnFound = func(imf *inMemoryFile) (*inMemoryFile, error) {
		return imf, nil
	}
	// check if file is a directory, otherwise return an error
	checkIfDirectory = func(imf *inMemoryFile) (*inMemoryFile, error) {
		if imf.info.fileType == file.Directory {
			return imf, nil
		}
		return nil, fserrors.ErrInvalidFileType
	}
)

// navigateToLastDirInPath navigates to the last Dir returned by filepath.Dir() in the path following symbolic links.
// If specied it creates any missing intermediate directories.
// If symbolic link points to an invalid location, it returns an error even if createMissingDir is true.
func (fs *MemoryFileSystem) navigateToLastDirInPath(path *fspath.FileSystemPath, workingDir file.File, createMissingDir bool, linkDepth int) (*inMemoryFile, error) {
	// Find path starting point
	pathRoot, err := fs.findPathRoot(path, workingDir)
	if err != nil {
		return nil, err
	}

	// Move through all the directories in the path, create missing if needed
	pathDirs := pathDirs(path, workingDir)
	if createMissingDir {
		return fs.lookupDirAndCreateMissingDirectories(pathRoot, pathDirs, linkDepth)
	}
	return fs.lookupDir(pathRoot, pathDirs, linkDepth)
}

// navigateToEndOfPath navigates to the last directory/file in the path following symbolic links.
// If specied it creates any missing intermediate directory.
// If symbolic link points to an invalid location, it returns an error even if createMissingDir is true.
func (fs *MemoryFileSystem) navigateToEndOfPath(path *fspath.FileSystemPath, workingDir file.File, createMissingDir bool, linkDepth int) (*inMemoryFile, error) {
	// Find path starting point
	pathRoot, err := fs.findPathRoot(path, workingDir)
	if err != nil {
		return nil, err
	}

	// Find where to add the file, and create intermediate directories if needed
	pathDirs := pathNames(path, workingDir)
	if createMissingDir {
		return fs.lookupFileAndCreateMissingDirectories(pathRoot, pathDirs, linkDepth)
	}
	return fs.lookupFile(pathRoot, pathDirs, linkDepth)
}

// appendMatchingFiles walks the file system and appends any file matching the specified name.
// if current file is a directory, recursively append every matching file in the subtree.
func (fs *MemoryFileSystem) appendMatchingFiles(matchingFiles []file.FileInfo, dir *inMemoryFile, name string) ([]file.FileInfo, error) {
	err := fs.walk(dir, func(fileName string, imf *inMemoryFile) error {
		var err error
		// add matching file
		if fileName == name {
			matchingFiles = append(matchingFiles, imf.Info())
		}

		// if directory, go down the tree
		if imf.info.fileType == file.Directory {
			matchingFiles, err = fs.appendMatchingFiles(matchingFiles, imf, name)
		}
		return err
	})

	return matchingFiles, err
}

func (fs *MemoryFileSystem) lookupDirAndCreateMissingDirectories(pathRoot *inMemoryFile, pathNames []string, linkDepth int) (*inMemoryFile, error) {
	return fs.doLookupFile(pathRoot, pathNames, checkIfDirectory, fs.createDirectory, linkDepth)
}

func (fs *MemoryFileSystem) lookupDir(pathRoot *inMemoryFile, pathNames []string, linkDepth int) (*inMemoryFile, error) {
	return fs.doLookupFile(pathRoot, pathNames, checkIfDirectory, errorOnNotFoundFn, linkDepth)
}

func (fs *MemoryFileSystem) lookupFileAndCreateMissingDirectories(pathRoot *inMemoryFile, pathNames []string, linkDepth int) (*inMemoryFile, error) {
	return fs.doLookupFile(pathRoot, pathNames, noopOnFound, fs.createDirectory, linkDepth)
}

func (fs *MemoryFileSystem) lookupFile(pathRoot *inMemoryFile, pathNames []string, linkDepth int) (*inMemoryFile, error) {
	return fs.doLookupFile(pathRoot, pathNames, noopOnFound, errorOnNotFoundFn, linkDepth)
}

// doLookupFile moves through every directory/file in pathnames untile it reaches the end of the array or
// an error occurs
func (fs *MemoryFileSystem) doLookupFile(pathRoot *inMemoryFile, pathNames []string, onFound onFileFoundFn, onNotFound onFileNotFoundFn, linkDepth int) (*inMemoryFile, error) {
	currentFile := pathRoot
	for _, nextFileName := range pathNames {
		var err error
		// move to next node in the file tree
		currentFile, err = fs.moveToNextFile(currentFile, nextFileName, onFound, onNotFound, linkDepth)
		if err != nil {
			return nil, err
		}
	}

	return currentFile, nil
}

// moveToNextFile moves to the next file in the path.
// If the file is not found the "onNotFound" callback is called.
// If file is a symlink, the symlink is resolved to the original file.
// If file is found, the "onFound" callback is called.
func (fs *MemoryFileSystem) moveToNextFile(currentFile *inMemoryFile, nextFileName string, onFound onFileFoundFn, onNotFound onFileNotFoundFn, linkDepth int) (*inMemoryFile, error) {
	nextFile, found := currentFile.fileMap[nextFileName]
	if !found {
		return onNotFound(nextFileName, currentFile)
	}

	if nextFile.info.fileType != file.SymbolicLink {
		return onFound(nextFile)
	}

	return fs.resolveSymlink(nextFile, linkDepth)
}

// resolveSymlink tries to resolve symlink and returns an error if the link points to a file
// that does not exixt or too many symlink were followed.
func (fs *MemoryFileSystem) resolveSymlink(currentFile *inMemoryFile, linkDepth int) (*inMemoryFile, error) {
	if linkDepth >= MAX_LINK_DEPTH {
		return nil, fserrors.ErrTooManyLinks
	}

	// Link contains absolute path, so no need to pass working dir
	return fs.navigateToEndOfPath(fspath.NewFileSystemPath(currentFile.link), nil, false, linkDepth+1)
}

// findPathRoot finds the path starting point
func (fs *MemoryFileSystem) findPathRoot(path *fspath.FileSystemPath, workingDir file.File) (*inMemoryFile, error) {
	if path.IsAbs() || workingDir == nil {
		return fs.root, nil
	}

	return fs.resolveWorkDir(path, workingDir)
}
