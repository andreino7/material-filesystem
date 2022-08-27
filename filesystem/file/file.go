package file

type FileType int

const (
	RegularFile FileType = iota
	Directory
	SymbolicLink
)

type Permission int

const (
	RW Permission = iota
	RO
	NONE
)

type FileInfo interface {
	Name() string
	FileType() FileType
	AbsolutePath() string
	UserId() string
	GroupId() string
}

type FileData interface {
	Data() []byte
}

type FilePermissions interface {
	World() Permission
	User() Permission
	Group() Permission
}

type File interface {
	Info() FileInfo
	Data() FileData
	Permissions() FilePermissions
}
