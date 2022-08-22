package memoryfs

import (
	"fmt"
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fspath"
	"sort"
)

type onFileNotFoundFn func(string, *inMemoryFile) (*inMemoryFile, error)
type onFileFoundFn func(*inMemoryFile) (*inMemoryFile, error)

var (
	errorOnNotFoundFn = func(filename string, parent *inMemoryFile) (*inMemoryFile, error) {
		return nil, fmt.Errorf("no such file or directory")
	}
	noopOnFound = func(file *inMemoryFile) (*inMemoryFile, error) {
		return file, nil
	}
	checkIfDirectory = func(file *inMemoryFile) (*inMemoryFile, error) {
		if file.Info().IsDirectory() {
			return file, nil
		}
		return nil, fmt.Errorf("file is not a directory")
	}
)

func (fs *MemoryFileSystem) FindFiles(name string, path *fspath.FileSystemPath, workingDir file.File) ([]file.FileInfo, error) {
	// Initialize result
	matchingFiles := []file.FileInfo{}

	if err := checkFileName(name); err != nil {
		return nil, err
	}

	fs.mutex.RLock()
	defer fs.mutex.RUnlock()

	// Get directory to start the search
	dir, err := fs.GetDirectory(path, workingDir)
	if err != nil {
		return matchingFiles, err
	}

	// this cast is safe because GetDirectory always returns "inMemoryFile"
	inMemoryDir := dir.(*inMemoryFile)
	matchingFiles = fs.appendMatchingFiles(matchingFiles, inMemoryDir, name)

	// sort lexicographically
	sort.Sort(ByAbsolutePath(matchingFiles))
	return matchingFiles, nil
}

func (fs *MemoryFileSystem) navigateToLastDirInPath(path *fspath.FileSystemPath, workingDir file.File, createMissingDir bool) (*inMemoryFile, error) {
	// Find path starting point
	pathRoot, err := fs.findPathRoot(path, workingDir)
	if err != nil {
		return nil, err
	}

	// Move through all the directories in the path, create missing if needed
	pathDirs := pathDirs(path, workingDir)
	if createMissingDir {
		return fs.lookupDirAndCreateMissingDirectories(pathRoot, pathDirs)
	}
	return fs.lookupDir(pathRoot, pathDirs)
}

func (fs *MemoryFileSystem) navigateToEndOfPath(path *fspath.FileSystemPath, workingDir file.File, createMissingDir bool) (*inMemoryFile, error) {
	// Find path starting point
	pathRoot, err := fs.findPathRoot(path, workingDir)
	if err != nil {
		return nil, err
	}

	// Find where to add the file, and create intermediate directories if needed
	pathDirs := pathNames(path, workingDir)
	if createMissingDir {
		return fs.lookupFileAndCreateMissingDirectories(pathRoot, pathDirs)
	}
	return fs.lookupFile(pathRoot, pathDirs)
}

func (fs *MemoryFileSystem) appendMatchingFiles(matchingFiles []file.FileInfo, dir *inMemoryFile, name string) []file.FileInfo {
	for fileName, file := range dir.fileMap {
		// skip special keys to avoid infinite cycle
		if fileName == ".." || fileName == "." || fileName == "/" {
			continue
		}

		// add matching file
		if fileName == name {
			matchingFiles = append(matchingFiles, file.Info())
		}

		// if directory, go down the tree
		if file.info.IsDirectory() {
			matchingFiles = fs.appendMatchingFiles(matchingFiles, file, name)
		}
	}
	return matchingFiles
}

func (fs *MemoryFileSystem) lookupDirAndCreateMissingDirectories(pathRoot *inMemoryFile, pathNames []string) (*inMemoryFile, error) {
	return fs.doLookupFile(pathRoot, pathNames, checkIfDirectory, fs.createDirectory)
}

func (fs *MemoryFileSystem) lookupDir(pathRoot *inMemoryFile, pathNames []string) (*inMemoryFile, error) {
	return fs.doLookupFile(pathRoot, pathNames, checkIfDirectory, errorOnNotFoundFn)
}

func (fs *MemoryFileSystem) lookupFileAndCreateMissingDirectories(pathRoot *inMemoryFile, pathNames []string) (*inMemoryFile, error) {
	return fs.doLookupFile(pathRoot, pathNames, noopOnFound, fs.createDirectory)
}

func (fs *MemoryFileSystem) lookupFile(pathRoot *inMemoryFile, pathNames []string) (*inMemoryFile, error) {
	return fs.doLookupFile(pathRoot, pathNames, noopOnFound, errorOnNotFoundFn)
}

func (fs *MemoryFileSystem) doLookupFile(pathRoot *inMemoryFile, pathNames []string, onFound onFileFoundFn, onNotFound onFileNotFoundFn) (*inMemoryFile, error) {
	currentFile := pathRoot
	for _, nextFileName := range pathNames {
		var err error
		// move to next node in the file tree
		currentFile, err = fs.moveToNextFile(currentFile, nextFileName, onFound, onNotFound)
		if err != nil {
			return nil, err
		}
	}

	return currentFile, nil
}

func (fs *MemoryFileSystem) moveToNextFile(currentFile *inMemoryFile, nextFileName string, onFound onFileFoundFn, onNotFound onFileNotFoundFn) (*inMemoryFile, error) {
	nextFile, found := currentFile.fileMap[nextFileName]
	if !found {
		return onNotFound(nextFileName, currentFile)
	}

	return onFound(nextFile)
}

func (fs *MemoryFileSystem) findPathRoot(path *fspath.FileSystemPath, workingDir file.File) (*inMemoryFile, error) {
	if path.IsAbs() || workingDir == nil {
		return fs.root, nil
	}

	return fs.resolveWorkDir(path, workingDir)
}
