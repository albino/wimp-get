package main

import (
	"os"
	"fmt"
	"regexp"
	"encoding/json"
	"path"
	"io/ioutil"
	"wimp-get/wimp"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <wimp id or url>\n", os.Args[0])
		os.Exit(-1)
	}

	var id string;
	wimpUrl, _ := regexp.Compile("^https?://play\\.wimpmusic\\.com/album/")
	if wimpUrl.MatchString(os.Args[1]) {
		id = wimpUrl.ReplaceAllString(os.Args[1], "")
	} else {
		id = os.Args[1]
	}

	exe, _ := os.Executable()
	wDir := path.Dir(exe)
	magicData, e := ioutil.ReadFile(wDir + "/magic.json")
	if e != nil {
		panic(e)
	}

	var magic map[string]interface{}
	e = json.Unmarshal(magicData, &magic)
	if e != nil {
		panic(e)
	}

	println("Looking up album...")

	album, e := wimp.GetAlbum(id, magic["sessionId"].(string))
	if e != nil {
		panic(e)
	}

	fmt.Printf("[ %s - %s (%d) ]\n", album.Artist, album.Title, album.Year)
}
