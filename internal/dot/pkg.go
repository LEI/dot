package dot

// Pkg task
type Pkg struct {
	Task
	Name   string
	Args   []string
	OS     []string
	Action string
	Type   string
}

func (p *Pkg) String() string {
	// return fmt.Sprintf("%s %s %s %s %s", p.Type, p.Action, p.Name, p.Args, p.OS)
	return p.Name // fmt.Sprintf("%s %s", p.Name, p.Args)
}
