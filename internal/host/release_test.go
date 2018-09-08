package host

import (
	"reflect"
	"sort"
	"testing"

	ini "gopkg.in/go-ini/ini.v1"
)

// /etc/lsb-release

var centos7 = []byte(`
NAME="CentOS Linux"
VERSION="7 (Core)"
ID="centos"
ID_LIKE="rhel fedora"
VERSION_ID="7"
PRETTY_NAME="CentOS Linux 7 (Core)"
`)

// /etc/os-release

var debian8 = []byte(`
PRETTY_NAME="Debian GNU/Linux 8 (jessie)"
NAME="Debian GNU/Linux"
VERSION_ID="8"
VERSION="8 (jessie)"
ID=debian
HOME_URL="http://www.debian.org/"
SUPPORT_URL="http://www.debian.org/support"
BUG_REPORT_URL="https://bugs.debian.org/"
`)

var debian9 = []byte(`
PRETTY_NAME="Debian GNU/Linux 9 (stretch)"
NAME="Debian GNU/Linux"
VERSION_ID="9"
VERSION="9 (stretch)"
ID=debian
HOME_URL="https://www.debian.org/"
SUPPORT_URL="https://www.debian.org/support"
BUG_REPORT_URL="https://bugs.debian.org/"
`)

var ubuntu14 = []byte(`
NAME="Ubuntu"
VERSION="14.04.5 LTS, Trusty Tahr"
ID=ubuntu
ID_LIKE=debian
PRETTY_NAME="Ubuntu 14.04.5 LTS"
VERSION_ID="14.04"

DISTRIB_ID=Ubuntu
DISTRIB_RELEASE=14.04
DISTRIB_CODENAME=trusty
DISTRIB_DESCRIPTION="Ubuntu 14.04.5 LTS"
`)

func TestParseRelease(t *testing.T) {
	releaseTests := []struct {
		in  []byte
		out []string
		// r   Release
	}{
		{centos7, []string{"centos", "centos7", "Core", "rhel", "fedora"}},
		{debian8, []string{"debian", "debian8", "jessie"}},
		{debian9, []string{"debian", "debian9", "stretch"}},
		{ubuntu14, []string{"ubuntu", "ubuntu14", "trusty", "debian"}},
	}
	for _, tt := range releaseTests {
		r := Release{}
		// ini.MapTo(&r, path)
		// bytes.NewReader(tt.in)
		if err := ini.MapTo(&r, tt.in); err != nil {
			t.Fatalf("could not parse ini: %s", err)
		}
		s := r.Parse()
		sort.Strings(s)
		sort.Strings(tt.out)
		if !reflect.DeepEqual(s, tt.out) {
			t.Fatalf("%v should be %v with %+v", s, tt.out, string(tt.in))
		}
	}
}
