package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"
)

// TODO: create missing directories is an option
// TODO: handle keeping the file open for one single writer and multipe readers
func (fs *MemoryFileSystem) AppendToFile(path *fspath.FileSystemPath, content []byte, workingDir file.File) error {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	parentDir, err := fs.navigateToLastDirInPath(path, workingDir, true, 0)
	if err != nil {
		return err
	}

	fileToWrite, err := fs.getFileToWrite(path, parentDir)
	if err != nil {
		return err
	}

	fileToWrite.data.data = append(fileToWrite.data.data, content...)
	return nil
}

func (fs *MemoryFileSystem) getFileToWrite(path *fspath.FileSystemPath, parentDir *inMemoryFile) (*inMemoryFile, error) {
	if parentDir.info.fileType != file.Directory {
		return nil, fserrors.ErrInvalidFileType
	}

	fileToWrite, found := parentDir.fileMap[path.Base()]
	if !found {
		return fs.createFile(path.Base(), file.RegularFile, parentDir)
	}

	if fileToWrite.info.FileType() != file.RegularFile {
		return nil, fserrors.ErrInvalidFileType
	}

	return fileToWrite, nil
}
