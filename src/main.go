package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
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
	flags := get_flags()
	rss_file := get_rss_file(flags)
	episodes := parse_rss_file(rss_file, flags)
	var wait sync.WaitGroup

	for _, episode := range episodes {
		wait.Add(1)
		go download_wrapper(episode, &wait)
	}
	add_date_file(rss_file)
	wait.Wait()
}

// Sometimes the rss files can give conflicting data so I need a way to manage what information comes from which fields.
func modify_rss_file_data(rss_file_content []byte, flags map[string]bool) []byte {
	rss_file_string := string(rss_file_content)
	if flags["itunes-title"] {
		rss_file_string = strings.Replace(rss_file_string, "<title>", "<junktitle>", -1)
		rss_file_string = strings.Replace(rss_file_string, "itunes:title", "title", -1)
	} else {
		rss_file_string = strings.Replace(rss_file_string, "itunes:title", "junktitle", -1)
	}

	if flags["itunes-season"] {
		rss_file_string = strings.Replace(rss_file_string, "<season>", "<junkseason>", -1)
		rss_file_string = strings.Replace(rss_file_string, "itunes:season", "season", -1)
	} else {
		rss_file_string = strings.Replace(rss_file_string, "itunes:season", "junkseason", -1)
	}

	modified_rss_file := []byte(rss_file_string)

	return modified_rss_file
}

func parse_rss_file(rss_file string, flags map[string]bool) []Episode {
	feed := rss_feed{}
	rss_file_content, err := ioutil.ReadFile(rss_file)
	if err != nil {
		panic(err)
	}
	modified_rss_feed := modify_rss_file_data(rss_file_content, flags)

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

// Some basic error handeling and go routine management in a wrapper for download_episode()
func download_wrapper(episode Episode, wait *sync.WaitGroup) {
	defer wait.Done()
	retries := 5
	iters := 0
	for iters <= retries {
		err := download_episode(episode)
		if err != nil {
			iters += 1
		} else {
			return
		}
	}
	fmt.Println("Failed to download %s", episode.title)
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
	if err != nil {
		return err
	}
	return nil
}

func verify_directory(season int) {
	dir := fmt.Sprintf("season%d", season)
	os.Mkdir(dir, 0755)
}

func add_date_file(rss_file string) {
	date := time.Now()
	date_file := "README"
	message := fmt.Sprintf("This archive created from %s at %v", rss_file, date)
	ioutil.WriteFile(date_file, []byte(message), 0644)
}

func download_rss_file(feed_url string) (string, error) {
	fmt.Println(6)
	fmt.Println(feed_url)
	resp, err := http.Get(feed_url)
	if err != nil {
		fmt.Println(2)
		return "", err
	}
	fmt.Println(7)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Could not get rss feed file")
	}
	fmt.Println(8)
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(3)
		return "", err
	}
	fmt.Println(9)
	out_file := strings.SplitAfter(feed_url, "/")[len(strings.SplitAfter(feed_url, "/"))-1]
	err = ioutil.WriteFile(out_file, content, 0644)
	if err != nil {
		fmt.Println(4)
		return "", err
	}
	fmt.Println(5)
	return out_file, nil

}
