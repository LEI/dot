package dot

// Hooks task list
type Hooks []*Hook

// Hook command to execute
type Hook struct {
	Command string
	Shell   string
	Action  string // install, remove
	OS      []string
}

// // InstallString string
// func (h *Hook) InstallString() string {
// 	return fmt.Sprintf("%s", h.Command)
// }

// // Install ...
// func (h *Hook) Install() error {
// 	fmt.Println("$", h.InstallString())
// 	return nil
// }

// // RemoveString string
// func (h *Hook) RemoveString() string {
// 	return fmt.Sprintf("%s", h.Command)
// }

// // Remove ...
// func (h *Hook) Remove() error {
// 	fmt.Println("$", h.RemoveString())
// 	return nil
// }
