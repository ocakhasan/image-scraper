package main

import (
	"flag"
	"log"

	"github.com/ocakhasan/image-scraper/scraper"
)

var folderFlag = flag.String("f", "", "folder for the images")
var websiteFlag = flag.String("w", "", "website to crawl images")

func main() {
	flag.Parse()
	err := scraper.GetImages(*websiteFlag, *folderFlag)
	if err != nil {
		log.Fatal(err)
	}
}
