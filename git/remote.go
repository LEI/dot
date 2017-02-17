package git

import (
	"fmt"
	"strings"
)

var (
	User  = "git"
	Host  = "github.com"
	Https bool // = true
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

func (r *Remote) SetUrl(url string) *Remote {
	switch {
	case strings.HasPrefix(url, "git@"),
		strings.HasPrefix(url, "git://"),
		// strings.HasPrefix(url, "http://"),
		strings.HasPrefix(url, "https://"),
		strings.HasPrefix(url, "ssh://"):
		r.URL = url
	default:
		if Https {
			r.URL = fmt.Sprintf("https://%s/%s", Host, url)
		} else {
			r.URL = fmt.Sprintf("%s@%s:%s", User, Host, url)
		}
	}
	if !strings.HasSuffix(url, ".git") {
		r.URL += ".git"
	}
	return r
}
