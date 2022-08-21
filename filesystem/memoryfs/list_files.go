package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fspath"
)

func (fs *MemoryFileSystem) ListFiles(path *fspath.FileSystemPath, workingDir file.File) ([]file.FileInfo, error) {
	fs.mutex.RLock()
	defer fs.mutex.RUnlock()
	return nil, nil

	// 	// Initialize result
	// 	files := []file.FileInfo{}

	// 	// lookup parent dir
	// 	parent, err := fs.lookupPathEndWithCreateMissingDir(path, workingDir, false)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	for name, file := range parent.fileMap {
	// 		// skip special entries
	// 		if name != ".." && name != "." && name != "/" {
	// 			files = append(files, file.Info())
	// 		}
	// 	}

	//		return files, nil
	//	}
}
