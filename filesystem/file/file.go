package file

type FileInfo struct {
	name        string
	isDirectory bool
}

type FileData struct {
}

type File struct {
	info *FileInfo
	data *FileData
}

func (f *File) Info() *FileInfo {
	return f.info
}

func (f *File) Data() *FileData {
	return f.data
}

func (info *FileInfo) IsDirectory() bool {
	return info.isDirectory
}

func (info *FileInfo) Name() string {
	return info.name
}

func NewFile(name string, isDirectory bool) *File {
	return &File{
		info: &FileInfo{
			name:        name,
			isDirectory: isDirectory,
		},
	}
}
