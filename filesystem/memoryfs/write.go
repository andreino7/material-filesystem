package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"
)

// TODO: create missing directories is an option
// TODO: handle keeping the file open for one single writer and multipe readers
func (fs *MemoryFileSystem) AppendAll(path *fspath.FileSystemPath, content []byte) error {
	fs.mutex.Lock()

	fileToWrite, err := fs.traverseToBaseAndCreateIntermediateDirs(path)
	if err != nil && err == fserrors.ErrNotExist {
		fs.mutex.Unlock()
		return err
	}
	if fileToWrite.info.fileType != file.RegularFile {
		fs.mutex.Unlock()
		return fserrors.ErrInvalidFileType
	}

	fileToWrite.data.mutex.Lock()
	defer fileToWrite.data.mutex.Unlock()
	fs.mutex.Unlock()

	fileToWrite.data.append(content)
	return nil
}

func (fs *MemoryFileSystem) WriteAt(fileDescriptor string, content []byte, pos int) (int, error) {
	// Read lock the open file table
	fs.openFiles.mutex.RLock()

	data, found := fs.openFiles.table[fileDescriptor]
	if !found {
		fs.openFiles.mutex.RUnlock()
		return 0, fserrors.ErrNotOpen
	}

	// Write lock file
	data.data.mutex.Lock()
	defer data.data.mutex.Unlock()

	// Unlock file table
	fs.openFiles.mutex.RUnlock()

	return data.data.writeAt(content, pos), nil
}
