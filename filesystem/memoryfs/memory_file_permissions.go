package memoryfs

import (
	"material/filesystem/filesystem/file"
)

type inMemoryFilePermissions struct {
	user  file.Permission
	group file.Permission
	world file.Permission
}

func (f *inMemoryFilePermissions) World() file.Permission {
	return f.world
}

func (f *inMemoryFilePermissions) Group() file.Permission {
	return f.group
}

func (f *inMemoryFilePermissions) User() file.Permission {
	return f.user
}
