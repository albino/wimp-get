// +build !windows

package platform

import(
	"path"
)

func SanitiseFilename(filename string) (newName string, e error) {
	newName = filename
	return
}

func DirOf(filename string) (dirname string, e error) {
	dirname = path.Dir(filename)
	return
}
