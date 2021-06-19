package main

import (
	"flag"
	"fmt"
	"github.com/ocakhasan/image-scraper/utils"
)

var folderFlag = flag.String("f", "imageFolder", "folder for the images")
var websiteFlag = flag.String("w", "https://google.com", "website to crawl images")

func main(){
	flag.Parse()
	err := utils.GetImages(*websiteFlag, *folderFlag)
	if err != nil {
		fmt.Printf("err: %v", err)
	}
}
