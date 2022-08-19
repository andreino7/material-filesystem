package memoryfs

import "material/filesystem/filesystem/file"

type fileWrapper struct {
	file     *file.File
	children map[string]*fileWrapper
}

func newFileWrapper(name string, isDirectory bool) *fileWrapper {
	file := file.NewFile(name, isDirectory)
	if isDirectory {
		return &fileWrapper{
			file:     file,
			children: nil,
		}
	}

	return &fileWrapper{
		file:     file,
		children: map[string]*fileWrapper{},
	}
}
