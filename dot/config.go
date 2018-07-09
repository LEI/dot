package dot

// Config ...
type Config struct {
	Roles []*Role
}

// Execute ...
func (c *Config) Execute() error {
	return nil
}
