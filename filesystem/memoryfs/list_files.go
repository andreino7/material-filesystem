package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fspath"
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
	for name, file := range dir.(*inMemoryFile).fileMap {
		// skip special entries
		if name != ".." && name != "." && name != "/" {
			files = append(files, file.Info())
		}
	}

	return files, nil
}
