package file

type FileType int

const (
	RegularFile FileType = iota
	Directory
	SymbolicLink
)

type FileInfo interface {
	Name() string
	FileType() FileType
	AbsolutePath() string
}

type FileData interface {
	Data() []byte
}

type File interface {
	Info() FileInfo
	Data() FileData
}
