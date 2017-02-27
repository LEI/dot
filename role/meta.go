package role

import (
	"fmt"
)

type Meta struct {
	Source string
	Target string
	Roles  []*Role
	// cfg   *config.Provider
}

func NewMeta() *Meta {
	return &Meta{Roles: make([]*Role, 0)}
}

func (m *Meta) String() string {
	return fmt.Sprintf("%s -> %s = %+v", m.Source, m.Target, m.Roles)
}
