package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

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
	images, err := extract(url)
	if err != nil {
		return err
	}

	if _, err := os.Stat(folder); err != nil {
		if os.IsNotExist(err) {
			// file does not exist
			fmt.Printf("Folder %s does not exist. Creating one...\n", folder)
			if err := os.Mkdir(folder, 0700); err != nil {
				return fmt.Errorf("error while creating folder %s : %v", folder, err)
			}
		} else {
			fmt.Printf("Folder %s already exists. Images will be downladed in that folder.\n", folder)
		}
	}

	for _, image := range images {
		err := getImageFromURl(image, folder)
		if err != nil {
			return fmt.Errorf("error in %s: %v", image, err)
		}
		fmt.Printf("Image %s is done\n", path.Base(image))
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

func getImageFromURl(url, folder string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return fmt.Errorf("got %d in %s", resp.StatusCode, url)
	}

	file, err := os.Create(filepath.Join(folder, path.Base(url)))
	defer file.Close()
	if err != nil {
		return err
	}
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
