package memoryfs

import (
	"material/filesystem/filesystem/user"
	"sync"
)

type userTable struct {
	users map[string]*inMemoryUser
	sync.RWMutex
}

type groupTable struct {
	groups map[string]bool
	sync.RWMutex
}

func newGroupTable(defaultGroup string) *groupTable {
	return &groupTable{
		groups: map[string]bool{defaultUser: true},
	}
}

func newUserTable(defaultUser string, defaultGroup string) *userTable {
	return &userTable{
		users: map[string]*inMemoryUser{defaultUser: newInMemoryUser(defaultUser, defaultGroup)},
	}
}

type inMemoryGroup struct {
	id string
}

type inMemoryUser struct {
	id           string
	groupSet     user.GroupSet
	primaryGroup string
}

func newInMemoryUser(id string, primaryGroup string) *inMemoryUser {
	return &inMemoryUser{
		id:           id,
		primaryGroup: primaryGroup,
		groupSet:     user.GroupSet{primaryGroup: true},
	}
}

func (u *inMemoryUser) Id() string {
	return u.id
}

func (u *inMemoryUser) Groups() user.GroupSet {
	return u.groupSet
}

func (g *inMemoryGroup) Id() string {
	return g.id
}
