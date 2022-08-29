package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fspath"
	"regexp"
	"sort"
)

func (fs *MemoryFileSystem) FindFiles(nameRegex string, path *fspath.FileSystemPath) ([]file.FileInfo, error) {
	// Initialize result
	matchingFiles := []file.FileInfo{}

	fs.RLock()
	defer fs.RUnlock()
	// Get directory to start the search
	exp, err := regexp.Compile(nameRegex)
	if err != nil {
		return matchingFiles, err
	}

	err = fs.Walk(path, func(f file.File) error {
		if exp.MatchString(f.Info().Name()) {
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
