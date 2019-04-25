package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type rss_feed struct {
	XMLName  xml.Name  `xml:"rss"`
	RSS_Show *RSS_Show `xml:"channel"`
}

type RSS_Show struct {
	XMLName      xml.Name      `xml:"channel"`
	title        string        `xml:"title"`
	RSS_Episodes []RSS_Episode `xml:"item"`
}

type RSS_Episode struct {
	XMLName       xml.Name      `xml":item"`
	Title         string        `xml:"title"`
	RSS_Enclosure RSS_Enclosure `xml:"enclosure"`
	Season        int           `xml:"season"`
}

type RSS_Enclosure struct {
	XMLName xml.Name `xml:"enclosure"`
	Url     string   `xml:"url,attr"`
}

type Episode struct {
	title  string
	url    string
	season int
}

func main() {
	rss_file := os.Args[1]
	episodes := parse_rss_file(rss_file)

	for _, episode := range episodes {
		go download_wrapper(episode)
	}
}

func parse_rss_file(rss_file string) []Episode {
	feed := rss_feed{}
	rss_file_content, err := ioutil.ReadFile(rss_file)
	if err != nil {
		panic(err)
	}

	modified_rss_feed := strings.Replace(string(rss_file_content), "itunes:season", "season", -1)

	decoder := xml.NewDecoder(bytes.NewReader([]byte(modified_rss_feed)))
	err = decoder.Decode(&feed)
	if err != nil {
		panic(err)
	}
	show := feed.RSS_Show
	episodes := show.RSS_Episodes
	episode_slice := []Episode{}

	for _, item := range episodes {
		this_episode := Episode{title: item.Title, url: item.RSS_Enclosure.Url, season: item.Season}
		episode_slice = append(episode_slice, this_episode)
	}

	return episode_slice
}

func download_wrapper(episode Episode) {
	err := download_episode(episode)
	if err != nil {
		fmt.Println(err)
	}
}

func download_episode(episode Episode) error {
	resp, err := http.Get(episode.url)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	out_file := fmt.Sprintf("season%d/%s.mp3", episode.season, episode.title)
	verify_directory(episode.season)
	err = ioutil.WriteFile(out_file, body, 0644)
	return nil
}

func verify_directory(season int) {
	dir := fmt.Sprintf("season%d", season)
	os.Mkdir(dir, 0755)
}
