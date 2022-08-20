package memoryfs

import (
	"fmt"
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fspath"
	"sort"
)

func (fs *MemoryFileSystem) FindFiles(name string, path *fspath.FileSystemPath, workingDir file.File) ([]file.FileInfo, error) {
	fs.mutex.RLock()
	defer fs.mutex.RUnlock()

	// Initialize result
	matchingFiles := []file.FileInfo{}

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

func (fs *MemoryFileSystem) lookupParentDirWithCreateMissingDir(path *fspath.FileSystemPath, workingDir file.File, createMissingDir bool) (*inMemoryFile, error) {
	// Find path starting point
	pathRoot, err := fs.findPathRoot(path, workingDir)
	if err != nil {
		return nil, err
	}

	// Find where to add the file, eventually create intermediate directories
	pathDirs := pathDirs(path, workingDir)
	return fs.lookupDirWithCreateMissing(pathRoot, pathDirs, createMissingDir)
}

func (fs *MemoryFileSystem) appendMatchingFiles(matchingFiles []file.FileInfo, dir *inMemoryFile, name string) []file.FileInfo {
	for fileName, file := range dir.fileMap {
		// skip special keys to avoid infinite cycle
		if fileName == ".." || fileName == "." {
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

func (fs *MemoryFileSystem) lookupDir(pathRoot *inMemoryFile, pathNames []string) (*inMemoryFile, error) {
	return fs.lookupDirWithCreateMissing(pathRoot, pathNames, false)
}

// TODO: refactor
func (fs *MemoryFileSystem) lookupDirWithCreateMissing(pathRoot *inMemoryFile, pathNames []string, createMissing bool) (*inMemoryFile, error) {
	var err error
	for _, currentDir := range pathNames {
		tmp, found := pathRoot.fileMap[currentDir]
		if !found {
			if createMissing {
				tmp, err = fs.createFile(currentDir, true, pathRoot)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, fmt.Errorf("no such file or directory")
			}
		}

		if !tmp.info.IsDirectory() {
			return nil, fmt.Errorf("file is not a directory")
		}

		pathRoot = tmp
	}

	return pathRoot, nil
}

func (fs *MemoryFileSystem) findPathRoot(path *fspath.FileSystemPath, workingDir file.File) (*inMemoryFile, error) {
	if path.IsAbs() || workingDir == nil {
		return fs.root, nil
	}

	return fs.resolveWorkDir(path, workingDir)
}
