package main

import (
	"os"
	"fmt"
	"regexp"
	"encoding/json"
	"path"
	"io"
	"io/ioutil"
	"wimp-get/wimp"
	"wimp-get/platform"
	"os/exec"
	"net/http"
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

	// Determine whether we have more than one disc
	multidisc := false
	for _, track := range album.Tracks {
		if track.Volume > 1 {
			multidisc = true
			break
		}
	}

	dirName := album.Artist+" - "+album.Title+" ("+fmt.Sprintf("%d", album.Year)+") [WEB FLAC]"

	e = os.Mkdir(dirName, os.FileMode(0755))
	if e != nil {
		panic(e)
	}

	// Time to do the ripping!
	for _, track := range album.Tracks {
		num := fmt.Sprintf("%d", track.Number)
		if len(num) < 2 {
			num = "0"+num
		}

		fmt.Printf("[%d/%s] %s...", track.Volume, num, track.Title)

		var filename string
		if (multidisc) {
			filename = fmt.Sprintf("%s/Disc %d/%s - %s.flac", dirName, track.Volume, num, track.Title)
		} else {
			filename = fmt.Sprintf("%s/%s - %s.flac", dirName, num, track.Title)
		}

		filename, e = platform.SanitiseFilename(filename)
		if e != nil {
			panic(e)
		}

		// create disc dir if necessary
		if _, e = os.Stat(path.Dir(filename)); e != nil {
			if os.IsNotExist(e) {
				e = os.Mkdir(path.Dir(filename), os.FileMode(0755))
				if e != nil {
					panic(e)
				}
			} else {
				panic(e)
			}
		}

		resp, e := http.Get(track.Url)
		if e != nil {
			panic(e)
		}

		ffmpeg := exec.Command(magic["ffmpeg"].(string), "-i", "-", "-metadata", "title="+track.Title, "-metadata", "artist="+track.Artist,
			"-metadata", "album="+album.Title, "-metadata", "year="+fmt.Sprintf("%d", album.Year), "-metadata", "track="+fmt.Sprintf("%d", track.Number),
			"-metadata", "albumartist="+album.Artist, "-metadata", "discnumber="+fmt.Sprintf("%d", track.Volume), filename)

		stdin, e := ffmpeg.StdinPipe()
		if e != nil {
			panic(e)
		}

		e = ffmpeg.Start()
		if e != nil {
			panic(e)
		}

		_, e = io.Copy(stdin, resp.Body)
		if e != nil {
			panic(e)
		}

		resp.Body.Close()

		println(" Done!")
	}
}
