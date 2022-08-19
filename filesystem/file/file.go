package file

type FileInfo interface {
	Name() string
	IsDirectory() bool
}

type FileData interface {
}

type File interface {
	Info() FileInfo
	Data() FileData
}
