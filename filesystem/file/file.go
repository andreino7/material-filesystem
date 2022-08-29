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
	// Name returns the file name
	Name() string
	// FileType returns the filetype: RegularFile, Directory, SymbolicLink
	FileType() FileType
	// AbsolutePath returns the file absolute path
	AbsolutePath() string
}

type FileData interface {
	// Data returns the file content
	Data() []byte
	// Size returns the size in byte
	Size() int
}

type File interface {
	// Info returns the file attributes
	Info() FileInfo
	// Data returns the file data
	Data() FileData
}
