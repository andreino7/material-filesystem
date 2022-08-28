package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fspath"
	"sort"
)

// ListFiles lists the files at the specified path
// sorted alphabetically.
// This implementation is thead safe.
//
// Returns an error when:
// - the file name is invalid
// - the target path is not a directory
// - the target path path does not exist
func (fs *MemoryFileSystem) ListFiles(path *fspath.FileSystemPath) ([]file.FileInfo, error) {
	fs.RLock()
	defer fs.RUnlock()

	// Initialize result
	files := []file.FileInfo{}

	// Get directory to list files
	dir, err := fs.GetDirectory(path)
	if err != nil {
		return files, err
	}

	// List all files
	err = fs.visitDir(dir.(*inMemoryFile), func(_ string, curr *inMemoryFile) error {
		files = append(files, curr.Info())
		return nil
	})
	if err != nil {
		return files, err
	}

	sort.Sort(ByAbsolutePath(files))
	return files, nil
}
