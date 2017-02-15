package git

import (
	"strings"
)

var (
	BaseUrl = "https://github.com"
)

type Remote struct {
	Name string
	URL  string
}

func NewRemote(name string, url string) *Remote {
	r := &Remote{Name: name}
	r.SetUrl(url)
	return r
}

func (r *Remote) SetUrl(value string) *Remote {
	var url string
	switch {
	case strings.HasPrefix(value, "git@"),
		strings.HasPrefix(value, "git://"),
		// strings.HasPrefix(value, "http://"),
		strings.HasPrefix(value, "https://"),
		strings.HasPrefix(value, "ssh://"):
		url = value
	default:
		url = BaseUrl + "/" + value
	}
	if !strings.HasSuffix(url, ".git") {
		url += ".git"
	}
	return r
}
