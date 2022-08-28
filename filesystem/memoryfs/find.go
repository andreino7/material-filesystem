package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fspath"
	"sort"
)

func (fs *MemoryFileSystem) FindFiles(name string, path *fspath.FileSystemPath) ([]file.FileInfo, error) {
	// Initialize result
	matchingFiles := []file.FileInfo{}

	if err := checkFileName(name); err != nil {
		return nil, err
	}

	fs.RLock()
	defer fs.RUnlock()
	// Get directory to start the search
	err := fs.Walk(path, func(f file.File) error {
		if f.Info().Name() == name {
			matchingFiles = append(matchingFiles, f.Info())
		}
		return nil
	}, func(f file.File) bool {
		return true
	}, true)

	if err != nil {
		return matchingFiles, err
	}

	sort.Sort(ByAbsolutePath(matchingFiles))

	return matchingFiles, nil
}
