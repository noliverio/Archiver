package main

import "os"
import "fmt"
import "log"

func get_flags() map[string]bool {
	flags := make(map[string]bool)

	if len(os.Args) == 1 {
		help()
	}

	for _, flag := range os.Args[2:] {
		if flag[:2] != "--" {
			log.Fatalln(fmt.Sprintf("Fatal error: malformed flag: %s", flag))
		} else {
			flags[flag[:2]] = true
		}
	}

	return flags
}

func get_rss_file(flags map[string]bool) string {
	var rss_file string
	if flags["net_rss"] {
		fmt.Println("Net based rss xml files not currently supported")
		os.Exit(1)
	} else {
		rss_file = os.Args[1]
	}
	return rss_file
}

func help() {
	fmt.Println(`Archiver is a tool for downloading podcast archives based on a provided rss xml file. 

Usage:

	archiver <feed.xml> [flags]

Optional flags:
	--itunes_title		Title from iTunes section of rss file is used instead of standard title field.
	--itunes_season		Season from iTunes section of rss file is used instead of standard season field.
	--net_rss		RSS file is retrieved over the network.
`)
	os.Exit(0)
}
