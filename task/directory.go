package task

type Directory struct {
	// *File
	Path string
}

func (d *Directory) Get() bool {
	fmt.Println("Get dir", d)
}

func (d *Directory) Check() bool {
	fmt.Println("Check dir", d)
}

func (d *Directory) Sync() bool {
	fmt.Println("Sync dir", d)
}
