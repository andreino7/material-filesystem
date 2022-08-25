package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fspath"
	"sort"
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
	matchingFiles, err = fs.appendMatchingFiles(matchingFiles, inMemoryDir, name)
	if err != nil {
		return matchingFiles, err
	}

	// sort lexicographically
	sort.Sort(ByAbsolutePath(matchingFiles))
	return matchingFiles, nil
}
