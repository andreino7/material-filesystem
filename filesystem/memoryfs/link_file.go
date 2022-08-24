package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"
)

// TODO: validate file name
// TODO: make create intermediate directories configurable
func (fs *MemoryFileSystem) CreateHardLink(srcPath *fspath.FileSystemPath, destPath *fspath.FileSystemPath, workingDir file.File) (file.FileInfo, error) {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	// Locate file to link
	fileToLink, err := fs.navigateToEndOfPath(srcPath, workingDir, false, 0)
	if err != nil {
		return nil, err
	}

	// Only hard links to regular file supported
	if fileToLink.info.fileType != file.RegularFile {
		return nil, fserrors.ErrInvalidFileType
	}

	// Create an empty file
	hardLink, err := fs.addFileToFsLockFree(destPath, workingDir, file.RegularFile, false)
	if err != nil {
		return nil, err
	}

	// Point the file to the same underline data
	hardLink.data = fileToLink.data
	return hardLink.info, nil
}

// TODO: make create intermediate directories configurable
func (fs *MemoryFileSystem) CreateSymbolicLink(srcPath *fspath.FileSystemPath, destPath *fspath.FileSystemPath, workingDir file.File) (file.FileInfo, error) {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	// Locate file to link
	fileToLink, err := fs.navigateToEndOfPath(srcPath, workingDir, false, 0)
	if err != nil {
		return nil, err
	}

	// Create an empty file
	symLink, err := fs.addFileToFsLockFree(destPath, workingDir, file.RegularFile, false)
	if err != nil {
		return nil, err
	}

	// Point the file to the original file
	symLink.link = fileToLink.info.absolutePath
	return symLink.info, nil
}
