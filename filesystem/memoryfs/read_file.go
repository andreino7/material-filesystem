package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"
)

// TODO: handle read from location instead of all content
func (fs *MemoryFileSystem) ReadFile(path *fspath.FileSystemPath, workingDir file.File) ([]byte, error) {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	fileToRead, err := fs.navigateToEndOfPath(path, workingDir, false, 0)
	if err != nil {
		return nil, err
	}

	return fs.getFileData(fileToRead)
}

func (fs *MemoryFileSystem) getFileData(fileToRead *inMemoryFile) ([]byte, error) {
	if fileToRead.info.fileType != file.RegularFile {
		return nil, fserrors.ErrInvalidFileType
	}

	return fileToRead.data.data, nil
}
