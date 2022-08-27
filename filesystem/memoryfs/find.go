package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fspath"
	"material/filesystem/filesystem/user"
	"sort"
)

func (fs *MemoryFileSystem) FindFiles(name string, path *fspath.FileSystemPath, user user.User) ([]file.FileInfo, error) {
	// Initialize result
	matchingFiles := []file.FileInfo{}

	if err := checkFileName(name); err != nil {
		return nil, err
	}

	fs.RLock()
	defer fs.RUnlock()

	// Get directory to start the search
	dir, err := fs.GetDirectory(path, user)
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
