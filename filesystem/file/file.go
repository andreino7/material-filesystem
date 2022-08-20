package file

type FileInfo interface {
	Name() string
	IsDirectory() bool
	AbsolutePath() string
}

type FileData interface {
}

type File interface {
	Info() FileInfo
	Data() FileData
}
