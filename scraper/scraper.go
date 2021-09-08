package scraper

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

func extract(url string) ([]string, error) {
	//if url does not have https://, convert it to https://{url}
	if !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}

	var images []string
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "img" {
			for _, img := range n.Attr {
				if img.Key != "src" {
					continue
				}
				image, err := resp.Request.URL.Parse(img.Val)
				if err != nil {
					continue
				}
				images = append(images, image.String())
			}
		}
	}

	forEachNode(doc, visitNode, nil)
	return returnUniqueImages(images), nil
}

func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}

func GetImages(url, folder string) error {
	if strings.TrimSpace(url) == "" || strings.TrimSpace(folder) == "" {
		return fmt.Errorf("folder or website is not set. Set values with --f and --w flags respectively")
	}
	images, err := extract(url)
	if err != nil {
		return err
	}

	if _, err := os.Stat(folder); err != nil {
		if os.IsNotExist(err) {
			// file does not exist
			if err := os.Mkdir(folder, 0700); err != nil {
				return fmt.Errorf("error while creating folder %s : %v", folder, err)
			}
		}
	}

	fmt.Printf("There are %d images in %v\n", len(images), url)
	outputChannel := make(chan string, len(images))

	var wg sync.WaitGroup
	for _, image := range images {
		wg.Add(1)
		go func(image string) {
			getImageFromURl(image, folder, outputChannel, &wg)
		}(image)
	}

	wg.Wait()
	close(outputChannel)

	for output := range outputChannel {
		fmt.Print(output)
	}
	return nil
}

func returnUniqueImages(images []string) []string {
	uniqueMap := make(map[string]struct{})
	for _, v := range images {
		uniqueMap[v] = struct{}{}
	}

	uniqueSlice := make([]string, 0, len(uniqueMap))
	for key := range uniqueMap {
		uniqueSlice = append(uniqueSlice, key)
	}
	return uniqueSlice
}

func getImageFromURl(url, folder string, outputChannel chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	imageName := path.Base(url)
	resp, err := http.Get(url)
	if err != nil {
		outputChannel <- fmt.Sprintf("Error in %s : %s\n", imageName, err.Error())
		return
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		outputChannel <- fmt.Sprintf("got %d in %s", resp.StatusCode, url)
		return
	}

	file, err := os.Create(filepath.Join(folder, imageName))
	if err != nil {
		outputChannel <- fmt.Sprintf("Error in %s : %s\n", imageName, err.Error())
		return
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		outputChannel <- fmt.Sprintf("Error in %s : %s\n", imageName, err.Error())
		return
	}
	outputChannel <- fmt.Sprintf("- Image %s is downloaded\n", imageName)
}
