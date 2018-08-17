package dot

import "fmt"

// Links task list
type Links []*Link

// Link task
type Link struct {
	Source string
	Target string
}

func (c *Link) String() string {
	return fmt.Sprintf("ln -s %s %s", c.Source, c.Target)
}
