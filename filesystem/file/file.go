package file

type FileInfo interface {
	Name() string
	IsDirectory() bool
	AbsolutePath() string
}

type FileData interface {
	Data() []byte
}

type File interface {
	Info() FileInfo
	Data() FileData
}
