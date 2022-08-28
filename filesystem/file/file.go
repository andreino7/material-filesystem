package file

type FileType int
type WalkFn func(File) error
type FilterFn func(File) bool

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
