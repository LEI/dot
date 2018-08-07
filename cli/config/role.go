package config

// type Roles []*Role
// func (roles *Roles) list() { }

// Role structure
type Role struct {
	Name string
	URL string
	OS []string
	Deps []string `mapstructure:"dependencies"`
	// Copy map[string]string
	Link map[string]string
	// Template map[string]string
}
