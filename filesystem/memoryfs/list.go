package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fspath"
	"material/filesystem/filesystem/user"
	"sort"
)

func (fs *MemoryFileSystem) ListFiles(path *fspath.FileSystemPath, user user.User) ([]file.FileInfo, error) {
	fs.RLock()
	defer fs.RUnlock()

	// Initialize result
	files := []file.FileInfo{}

	// Get directory to list files
	dir, err := fs.GetDirectory(path, user)
	if err != nil {
		return files, err
	}

	// List all files
	err = fs.walk(dir.(*inMemoryFile), user, func(_ string, curr *inMemoryFile) error {
		files = append(files, curr.Info())
		return nil
	})
	if err != nil {
		return files, err
	}

	sort.Sort(ByAbsolutePath(files))
	return files, nil
}
