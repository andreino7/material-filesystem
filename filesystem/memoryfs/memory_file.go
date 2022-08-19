package memoryfs

import "material/filesystem/filesystem/file"

type inMemoryFileInfo struct {
	name        string
	isDirectory bool
}

type inMemoryFileData struct {
}

type inMemoryFile struct {
	info      file.FileInfo
	data      file.FileData
	isDeleted bool
	children  map[string]*inMemoryFile
}

func (f inMemoryFile) Info() file.FileInfo {
	return f.info
}

func (f inMemoryFile) Data() file.FileData {
	return f.data
}

func (info inMemoryFileInfo) IsDirectory() bool {
	return info.isDirectory
}

func (info inMemoryFileInfo) Name() string {
	return info.name
}

func newInMemoryFile(name string, isDirectory bool) *inMemoryFile {
	info := &inMemoryFileInfo{
		name,
		isDirectory,
	}

	if !isDirectory {
		return &inMemoryFile{
			info:     info,
			children: nil,
		}
	}

	file := &inMemoryFile{
		info:     info,
		children: map[string]*inMemoryFile{},
	}

	file.children["."] = file
	return file
}
