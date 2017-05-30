// +build !windows

package platform

func SanitiseFilename(filename string) (newName string, e error) {
	newName = filename
	return
}
