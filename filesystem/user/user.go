package user

type GroupSet map[string]bool

type Group interface {
	Id() string
}

type User interface {
	Id() string
	Groups() GroupSet
	PrimaryGroup() string
}
