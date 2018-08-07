package config

// Dot structure
type Dot struct {
	Roles []*Role
	config *Config
}

// type Roles []*Role
// func (roles *Roles) list() { }

var (
	// DotConfig value
	DotConfig *Dot
)
