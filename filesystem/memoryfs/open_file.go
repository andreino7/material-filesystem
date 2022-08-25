package memoryfs

import (
	"material/filesystem/filesystem/file"
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

func (fs *MemoryFileSystem) Open(path *fspath.FileSystemPath, workingDir file.File) (string, error) {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	fileToOpen, err := fs.navigateToEndOfPath(path, workingDir, false, 0)
	if err != nil {
		return "", err
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
