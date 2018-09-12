package git

import (
	"fmt"
	"strings"
)

// const (
// 	gitFormat   = "%s@%s:%s.git" // git://
// 	sshFormat   = "ssh://%s@%s/%s.git"
// 	httpsFormat = "https://%s::@%s/%s.git"
// )

// // Remote URL
// type Remote struct {
// 	Scheme string
// 	User   string
// 	Host   string
// 	// *Remote
// }

// // Parse repo url
// func (r *Remote) Parse(url string) string {
// 	// return ParseURL(r.Scheme, r.User, r.Host, url)
// 	o := url
// 	if url == "" || strings.Contains(url, "://") ||
// 		strings.HasPrefix(url, r.User+"@"+r.Host) {
// 		return url // ErrEmptyURL
// 	}
// 	// Count the number of slashes (user/repo sep)
// 	ns := bytes.Count([]byte(url), []byte{'/'})
// 	if ns > 1 {
// 		return url
// 	}
// 	if ns == 0 {
// 		fmt.Fprintf(Stderr, "%s: missing slash in repo remote url\n", url)
// 		// url = fmt.Sprintf("%s/%s", User, url)
// 	}
// 	// r := &RemoteConfig{}
// 	// if ns == 1 {
// 	// 	parts := strings.SplitN(url, "/", 2)
// 	// 	owner := parts[0] // Git user name
// 	// 	name := parts[1]  // Project name
// 	// }
// 	// url = NewGitRemote(user, host, url).String()
// 	// url = NewHTTPSRemote(user, host, url).String()
// 	// url = NewSSHRemote(user, host, url).String()
// 	switch r.Scheme {
// 	case "git":
// 		url = fmt.Sprintf(gitFormat, r.User, r.Host, url)
// 	case "ssh":
// 		url = fmt.Sprintf(sshFormat, r.User, r.Host, url)
// 	case "https":
// 		fallthrough
// 	default:
// 		url = fmt.Sprintf(httpsFormat, r.User, r.Host, url)
// 	}
// 	fmt.Printf("---\nParsed\n'%s'\nto\n'%s'\nuser: %s\nhost: %s\n---\n", o, url, r.User, r.Host)
// 	return url
// }

// // Remote ...
// type Remote struct {
// 	// Args() []string
// 	// Format() string
// 	// args   []interface{}
// 	Format string
// 	*URL
// }

// // URL ...
// type URL struct {
// 	User string
// 	Host string
// 	Repo string
// }

// // FormatString ...
// func (r *Remote) String() string {
// 	// for _, i := range r.URL.Args() {
// 	// 	args = append(args, i.(string))
// 	// }
// 	repo := r.URL.Repo
// 	// if !strings.HasSuffix(repo, ".git") {
// 	// 	repo += ".git"
// 	// }
// 	return fmt.Sprintf(r.Format, r.URL.User, r.URL.Host, repo)
// }

// // NewRemote URL
// func NewRemote(format, user, host, repo string) *Remote {
// 	if host == "" {
// 		host = Host
// 	}
// 	if user == "" {
// 		user = User
// 	}
// 	return &Remote{
// 		format,
// 		&URL{user, host, repo},
// 	}
// }

// // NewGitRemote URL
// func NewGitRemote(user, host, repo string) *Remote {
// 	return NewRemote(gitFormat, user, host, repo)
// }

// // NewHTTPSRemote URL
// func NewHTTPSRemote(user, host, repo string) *Remote {
// 	if user != "" {
// 		user = fmt.Sprintf("%s::@", user)
// 	}
// 	return NewRemote(httpsFormat, user, host, repo)
// }

// // NewSSHRemote URL
// func NewSSHRemote(user, host, repo string) *Remote {
// 	return NewRemote(sshFormat, user, host, repo)
// }

// // ParseURL prepends user and host
// func ParseURL(proto, user, host, url string) string {
// 	o := url
// 	// if url != "" && !strings.Contains(url, "https://") {
// 	// 	url = fmt.Sprintf(urlFormat, url)
// 	// } else if url == "" && strings.Contains(dir, "/") && string(dir[0]) != "/" && string(dir[0]) != "~" {
// 	// 	url = fmt.Sprintf(urlFormat, dir)
// 	// }
// 	if url == "" || strings.Contains(url, "://") || strings.HasPrefix(url, user+"@"+host) {
// 		return url // ErrEmptyURL
// 	}
// 	// if !strings.HasPrefix(url, "https://") {
// 	// 	url = fmt.Sprintf(httpsFormat, User, Host, url)
// 	// }
// 	// TODO check if is not already an URL

// 	// Count the number of slashes (user/repo sep)
// 	ns := bytes.Count([]byte(url), []byte{'/'})
// 	if ns > 1 {
// 		return url
// 	}
// 	if ns == 0 {
// 		fmt.Fprintf(Stderr, "%s: missing slash in repo remote url", url)
// 		// url = fmt.Sprintf("%s/%s", User, url)
// 	}
// 	// r := &RemoteConfig{}
// 	// if ns == 1 {
// 	// 	parts := strings.SplitN(url, "/", 2)
// 	// 	owner := parts[0] // Git user name
// 	// 	name := parts[1]  // Project name
// 	// }
// 	// url = NewGitRemote(user, host, url).String()
// 	url = NewHTTPSRemote(user, host, url).String()
// 	// url = NewSSHRemote(user, host, url).String()
// 	fmt.Printf("---\nParsed\n'%s'\nto\n'%s'\nuser: %s\nhost: %s---\n", o, url, user, host)
// 	return url
// }

// CheckRemote URL of an existing repository
func CheckRemote(dir, url string) error {
	args := []string{"-C", dir, "config", "--local", "--get", "remote.origin.url"}
	buf, err := gitCombined(args...)
	if err != nil {
		if err.Error() == "exit status 128" {
			// Not a git repository, or not yet cloned?
			return nil
		}
		return fmt.Errorf("%s: remote check failed: %s", dir, err)
	}
	actual := string(buf)
	// TODO: check domain and `user/repo`
	actual = parseRepo(actual)
	repo := parseRepo(url)
	if !strings.Contains(actual, repo) {
		// log.Fatalf()
		return fmt.Errorf("remote mismatch: url is '%s' but repo has '%s'", url, actual)
	}
	return nil
}

func parseRepo(str string) string {
	str = strings.TrimSuffix(str, ".git")
	str = strings.Replace(str, ":", "/", 1)
	parts := strings.Split(str, "/")
	if len(parts) > 1 {
		return strings.Join(parts[len(parts)-2:], "/")
	}
	return str
}
