package wimp

import (
	"encoding/json"
	"net/http"
	"io/ioutil"
	"strconv"
	"strings"
	"fmt"
)

type Track struct {
	Title string
	Artist string
	Url string
	Number int
	Volume int
}

type Album struct {
	Title string
	Artist string
	Year int
	CoverUrl string
	// TODO Genre?
	Tracks []Track
}

func apiRequest(path, sessionId string) (sData []byte, e error) {
	resp, e := http.Get("https://play.wimpmusic.com/v1/"+path+"&countryCode=GB&sessionId="+sessionId)
	if e != nil {
		return
	}

  sData, e = ioutil.ReadAll(resp.Body)
	return
}

func getTrackUrl(id, sessionId string) (url string, e error) {
	resp, e := apiRequest("tracks/"+id+"/streamUrl?soundQuality=LOSSLESS", sessionId);
	if e != nil {
		return
	}

	var trackUrlData map[string]interface{}
	e = json.Unmarshal(resp, &trackUrlData)
	if e != nil {
		return
	}

	url = trackUrlData["url"].(string)
	return
}

func getTracks(id, sessionId string) (tracks []Track, e error) {
	resp, e := apiRequest("albums/"+id+"/tracks?", sessionId)
	if e != nil {
		return
	}

	var tracksData map[string]interface{}
	e = json.Unmarshal(resp, &tracksData)
	if e != nil {
		return
	}

	uTracks := tracksData["items"].([]interface{})

	for _, el := range uTracks {
		tr := el.(map[string]interface{})

		track := Track{
			Title: tr["title"].(string),
			Number: int(tr["trackNumber"].(float64)),
			Volume: int(tr["volumeNumber"].(float64)),
		}

		artist := tr["artist"].(map[string]interface{})
		track.Artist = artist["name"].(string)

		url, e := getTrackUrl(fmt.Sprintf("%.0f", tr["id"].(float64)), sessionId)
		if e != nil {
			return tracks, e
		}
		track.Url = url

		tracks = append(tracks, track)
	}

	return
}

func GetAlbum(id, sessionId string) (album Album, e error) {
	resp, e := apiRequest("albums/"+id+"?", sessionId)
	if e != nil {
		return
	}

	var albumData map[string]interface{}
	e = json.Unmarshal(resp, &albumData)
	if e != nil {
		return
	}

	released := albumData["releaseDate"].(string)
	year, e := strconv.Atoi(released[:4])
	if e != nil {
		return
	}

	album.Year = year

	artists := albumData["artists"].([]interface{})
	artist := artists[0].(map[string]interface{})
	album.Artist = artist["name"].(string)

	album.Title = albumData["title"].(string)

	coverPath := strings.Replace(albumData["cover"].(string), "-", "/", -1)
	album.CoverUrl = "http://resources.wimpmusic.com/images/"+coverPath+"/1280x1280.jpg"

	tracks, e := getTracks(id, sessionId)
	if e != nil {
		return
	}

	album.Tracks = tracks

	return
}
