// +build windows

package platform

import(
	"regexp"
)

func DirOf(filename string) (dirname string, e error) {
	r, e := regexp.Compile(`\\[^\\]+$`)
	if e != nil {
		return
	}

	dirname = r.ReplaceAllString(filename, "")
	return
}
