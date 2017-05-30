// +build windows
// currently untested as I don't have a windows machine to test on

package platform

import(
	"regexp"
)

func SanitiseFilename(filename string) (newName string, e error) {
	r, e := regexp.Compile("[\\?<>:\"/\\\\|\\*]")
	if e != nil {
		return
	}

	newName = r.ReplaceAllString(filename, "")
	return
}
