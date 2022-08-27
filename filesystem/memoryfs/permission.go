package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/user"
)

type permissionSet map[file.Permission]bool

var (
	read  = map[file.Permission]bool{file.RO: true, file.RW: true}
	write = map[file.Permission]bool{file.RW: true}
)

// TODO: document assumptions
func checkReadPermissions(f *inMemoryFile, user user.User) error {
	if checkPermissions(f, user, read) {
		return nil
	}
	return fserrors.ErrPermissionDenied
}

// TODO: document assumptions
func checkWritePermission(f *inMemoryFile, user user.User) error {
	if checkPermissions(f, user, write) {
		return nil
	}
	return fserrors.ErrPermissionDenied
}

func checkPermissions(f *inMemoryFile, user user.User, permissionSet permissionSet) bool {
	if user.Id() == "root" {
		return true
	}

	// check world permissions
	if _, found := permissionSet[f.permissions.world]; found {
		return true
	}

	// check user permissions
	if checkGroupPermissions(f, user, permissionSet) {
		return true
	}

	return checkUserPermissions(f, user, permissionSet)
}

func checkUserPermissions(f *inMemoryFile, user user.User, permissionSet permissionSet) bool {
	if f.userId != user.Id() {
		return false
	}
	_, found := permissionSet[f.permissions.user]
	return found
}

func checkGroupPermissions(f *inMemoryFile, user user.User, permissionSet permissionSet) bool {
	// check if user belongs to file's group
	if _, found := user.Groups()[f.groupId]; !found {
		return false
	}
	_, found := permissionSet[f.permissions.group]
	return found
}
