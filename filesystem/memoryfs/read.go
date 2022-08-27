package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"
)

// TODO: handle read from location instead of all content
func (fs *MemoryFileSystem) ReadAll(path *fspath.FileSystemPath) ([]byte, error) {
	fs.Lock()
	defer fs.Unlock()

	fileToRead, err := fs.traverseToBase(path)
	if err != nil {
		return nil, err
	}

	return fs.getFileData(fileToRead)
}

func (fs *MemoryFileSystem) ReadAt(fileDescriptor string, startPos int, endPos int) ([]byte, error) {
	if err := validatePos(startPos, endPos); err != nil {
		return nil, err
	}

	// Read lock the open file table
	fs.openFiles.RLock()

	data, found := fs.openFiles.table[fileDescriptor]
	if !found {
		fs.openFiles.RUnlock()
		return nil, fserrors.ErrNotOpen
	}

	// Write lock file
	data.data.RLock()
	defer data.data.RUnlock()

	// Unlock file table
	fs.openFiles.RUnlock()

	return data.data.readAt(startPos, endPos), nil
}

func (fs *MemoryFileSystem) getFileData(fileToRead *inMemoryFile) ([]byte, error) {
	if fileToRead.info.fileType != file.RegularFile {
		return nil, fserrors.ErrInvalidFileType
	}

	return fileToRead.data.data, nil
}

func validatePos(start int, end int) error {
	if end < start {
		return fserrors.ErrInvalid
	}

	if start < 0 {
		return fserrors.ErrInvalid
	}

	return nil
}
