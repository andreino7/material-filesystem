package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"
	"material/filesystem/filesystem/user"
	"sync"

	"github.com/google/uuid"
)

type fileDataWrapper struct {
	data *inMemoryFileData
	pos  int
}

type fileTable struct {
	table map[string]*fileDataWrapper
	sync.RWMutex
}

func newFileTable() *fileTable {
	return &fileTable{
		table: map[string]*fileDataWrapper{},
	}
}

// Open opens the file at the given location and returns the file descriptor.
// Returns an error if the file was not found or it's not a "regular" file.
func (fs *MemoryFileSystem) Open(path *fspath.FileSystemPath, user user.User) (string, error) {
	fs.Lock()
	defer fs.Unlock()

	fileToOpen, err := fs.traverseToBase(path, user)
	if err != nil {
		return "", err
	}

	if fileToOpen.info.fileType != file.RegularFile {
		return "", fserrors.ErrInvalidFileType
	}

	fs.openFiles.Lock()
	defer fs.openFiles.Unlock()

	// TODO: fd could just be an int
	fd := uuid.NewString()
	// TODO set pos to size
	fs.openFiles.table[fd] = &fileDataWrapper{data: fileToOpen.data, pos: 0}
	return fd, nil
}

func (fs *MemoryFileSystem) Close(descriptor string, user user.User) {
	fs.openFiles.Lock()
	defer fs.openFiles.Unlock()
	delete(fs.openFiles.table, descriptor)
}
