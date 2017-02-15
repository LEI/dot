package git

type Remote struct {
	// Name string
	URL  string
}

func NewRemote(name string, url string) *Remote {
	return &Remote{
		// Name: name,
		URL: UrlScheme(url),
	}
}

func (r *Remote) SetUrl(url string) *Remote {
	switch {
	case strings.HasPrefix(url, "git@"),
		strings.HasPrefix(url, "git://"),
		strings.HasPrefix(url, "https://"),
		strings.HasPrefix(url, "ssh://"):
	default:
		url := "https://github.com/" + str + ".git"
	}
	// if !strings.EndsWith(str, ".git") {
	// 	url += ".git"
	// }
	return r
}
