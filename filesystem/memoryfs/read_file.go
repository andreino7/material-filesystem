package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"
)

// TODO: handle read from location instead of all content
func (fs *MemoryFileSystem) ReadAll(path *fspath.FileSystemPath, workingDir file.File) ([]byte, error) {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	fileToRead, err := fs.navigateToEndOfPath(path, workingDir, false, 0)
	if err != nil {
		return nil, err
	}

	return fs.getFileData(fileToRead)
}

func (fs *MemoryFileSystem) ReadAt(fileDescriptor string, startPos int, endPos int) ([]byte, error) {
	// Read lock the open file table
	fs.openFiles.mutex.RLock()

	data, found := fs.openFiles.table[fileDescriptor]
	if !found {
		fs.openFiles.mutex.RUnlock()
		return nil, fserrors.ErrNotOpen
	}

	// Write lock file
	data.data.mutex.RLock()
	defer data.data.mutex.RUnlock()

	// Unlock file table
	fs.openFiles.mutex.RUnlock()

	return data.data.readAt(startPos, endPos), nil
}

func (fs *MemoryFileSystem) getFileData(fileToRead *inMemoryFile) ([]byte, error) {
	if fileToRead.info.fileType != file.RegularFile {
		return nil, fserrors.ErrInvalidFileType
	}

	return fileToRead.data.data, nil
}
