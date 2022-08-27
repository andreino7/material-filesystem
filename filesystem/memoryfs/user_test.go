package memoryfs_test

import "material/filesystem/filesystem/user"

type testUser struct {
	id           string
	groups       user.GroupSet
	primaryGroup string
}

var rootUser = newTestUser("root", "root")

func newTestUser(id string, primaryGroup string) *testUser {
	return &testUser{
		id:           id,
		primaryGroup: primaryGroup,
	}
}

func (u *testUser) Id() string {
	return u.id
}

func (u *testUser) Groups() user.GroupSet {
	return u.groups
}

func (u *testUser) PrimaryGroup() string {
	return u.primaryGroup
}
