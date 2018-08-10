package tasks

// Template task
type Template struct {
	Source, Target string
	Env            map[string]string
	Vars           map[string]interface{}
	// backup bool
}
