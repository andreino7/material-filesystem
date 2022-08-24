package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fspath"
	"sort"
)

func (fs *MemoryFileSystem) ListFiles(path *fspath.FileSystemPath, workingDir file.File) ([]file.FileInfo, error) {
	fs.mutex.RLock()
	defer fs.mutex.RUnlock()

	// Initialize result
	files := []file.FileInfo{}

	// Get directory to list files
	dir, err := fs.GetDirectory(path, workingDir)
	if err != nil {
		return files, err
	}

	// List all files
	err = fs.walk(dir.(*inMemoryFile), func(_ string, curr *inMemoryFile) error {
		files = append(files, curr.Info())
		return nil
	})
	if err != nil {
		return files, err
	}

	sort.Sort(ByAbsolutePath(files))
	return files, nil
}
