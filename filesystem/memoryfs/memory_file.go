package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fspath"
)

type inMemoryFile struct {
	info        *inMemoryFileInfo
	data        *inMemoryFileData
	permissions *inMemoryFilePermissions
	isDeleted   bool
	fileMap     map[string]*inMemoryFile
	link        *fspath.FileSystemPath
	userId      string
	groupId     string
}

func (f *inMemoryFile) Info() file.FileInfo {
	return f.info
}

func (f *inMemoryFile) Data() file.FileData {
	return f.data
}

func (f *inMemoryFile) Permissions() file.FilePermissions {
	return f.permissions
}

func newInMemoryFile(absolutePath string, fileType file.FileType, userId string, groupId string) *inMemoryFile {
	info := &inMemoryFileInfo{
		absolutePath: absolutePath,
		fileType:     fileType,
		userId:       userId,
		groupId:      groupId,
	}

	// TODO: document default
	permissions := &inMemoryFilePermissions{
		user:  file.RW,
		group: file.RW,
		world: file.RW,
	}

	newFile := &inMemoryFile{
		info:        info,
		data:        &inMemoryFileData{},
		fileMap:     map[string]*inMemoryFile{},
		permissions: permissions,
	}

	if fileType == file.RegularFile {
		newFile.data = &inMemoryFileData{}
	}

	newFile.fileMap["."] = newFile
	return newFile
}
