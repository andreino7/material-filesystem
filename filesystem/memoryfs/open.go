package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"
	"sync"

	"github.com/google/uuid"
)

type fileDataWrapper struct {
	data *inMemoryFileData
	pos  int
}

type FileTable struct {
	table map[string]*fileDataWrapper
	mutex sync.RWMutex
}

func (fs *MemoryFileSystem) Open(path *fspath.FileSystemPath) (string, error) {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	fileToOpen, err := fs.traverseToBase(path)
	if err != nil {
		return "", err
	}

	if fileToOpen.info.fileType != file.RegularFile {
		return "", fserrors.ErrInvalidFileType
	}

	fs.openFiles.mutex.Lock()
	defer fs.openFiles.mutex.Unlock()

	fd := uuid.NewString()
	// TODO set pos to size
	fs.openFiles.table[fd] = &fileDataWrapper{data: fileToOpen.data, pos: 0}
	return fd, nil
}

func (fs *MemoryFileSystem) Close(descriptor string) {
	fs.openFiles.mutex.Lock()
	defer fs.openFiles.mutex.Unlock()
	delete(fs.openFiles.table, descriptor)
}
