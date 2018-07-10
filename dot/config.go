package dot

import (
	// "fmt"
)

// Config ...
type Config struct {
	Roles []*Role
}

// Execute ...
func (c *Config) Execute() error {
	// fmt.Println("Executing...")
	for _, r := range c.Roles {
		if err := r.Execute(); err != nil {
			return err
		}
	}
	return nil
}
