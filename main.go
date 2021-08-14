package main

import (
	"flag"
	"fmt"

	"github.com/ocakhasan/image-scraper/scraper"
)

var folderFlag = flag.String("f", "imageFolder", "folder for the images")
var websiteFlag = flag.String("w", "https://google.com", "website to crawl images")

func main() {
	flag.Parse()
	err := scraper.GetImages(*websiteFlag, *folderFlag)
	if err != nil {
		fmt.Printf("err: %v", err)
	}
}
