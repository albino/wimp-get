// +build !windows

package platform

import(
	"path"
)

func DirOf(filename string) (dirname string, e error) {
	dirname = path.Dir(filename)
	return
}
