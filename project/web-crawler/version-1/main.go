/*
	Remi Henriot: henriot.rem@gmail.com
	I didn't get the time to finalize the program in order to resume the previous crawl but most of the work is done.
	Remaining tasks would be :
		* To scan the cr.destinationPath in order to initialize the cr.status HashMap and the cr.stack Slice
		* Change the status values as global and constant variables
	To go further :
		* Create a distributed crawler working across different machines
	To start the program in command line exec : "./web-crawler https://www.golangr.com /Users/remihenriot/tmp"
*/

package main

import (
	"fmt"
	"os"
	"net/http"
	"io/ioutil"
	"regexp"
	"strings"
	"sync"
)

type Crawler struct {
	mu sync.Mutex
	status map[string]string
	stack []string
	availableWorkers int
	input chan string
	startUrl, destinationPath string
}

func (cr *Crawler) init() {
    cr.status = make(map[string]string)
    cr.status[cr.startUrl] = "ready"
	cr.stack = append(cr.stack, cr.startUrl)
    cr.input = make(chan string, cr.availableWorkers)
}

// Download the web page and return content
func downloadWebPage(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return content
}

// Save content of the webpage or rename existing file
func saveWebPage(cr *Crawler, url, extension string, content []byte) {
	name := strings.Replace(strings.Replace(url, ":", "*", -1), "/", "!", -1)
	filename := cr.destinationPath + "/" + name + "." + extension
	existing := cr.destinationPath + "/" + name
	if extension == "done" {
		existing += ".saved"
	} else {
		if extension == "saved" {
			existing += ".ready"
		}
		file, err := os.Create(existing)
		if err != nil {
			panic(err)
		}
		file.Write(content)
		file.Close()
	}
	if err := os.Rename(existing, filename); err != nil {
		panic(err)
	}
}

// Parse the web page to find child URLs
func parseWebPage(url, html string) []string {
	childUrls := []string{}
	re := regexp.MustCompile(`<a\s+(?:[^>]*?\s+)?href="([^"]*)"`)

	matches := re.FindAllStringSubmatch(string(html), -1)
	for _, match := range matches {
		url := match[1]
		if i := strings.IndexByte(url, '?'); i > -1 {
			url = url[:i]
		}
		childUrls = append(childUrls, url)
	}

	return childUrls
}

// Trim suffix
func trimSuffix(s, suffix string) string {
    if strings.HasSuffix(s, suffix) {
        s = s[:len(s)-len(suffix)]
    }
    return s
}

// Manage relative urls in webpage and remove unecessary backslash
func cleanUrls(parentUrl string, childUrls []string) []string {
	var cleanUrls []string
	parentUrl = trimSuffix(parentUrl, "/")
	for _, url := range childUrls {
		if len(url) > 0 && string(url[0]) == "/" {
			cleanUrls = append(cleanUrls, parentUrl + url)
		} else if strings.HasPrefix(url, parentUrl) {
			cleanUrls = append(cleanUrls, url)
		}
	}
	return cleanUrls
}

// Prepare child pages by saving a .ready file and adding the urls to the stack
func prepareChildWebPages(cr *Crawler, urls []string) {
	for _, url := range urls {
		cr.mu.Lock()
		_, ok := cr.status[url]
		cr.status[url] = "ready"
		cr.mu.Unlock()
		if !ok {
			saveWebPage(cr, url, "ready", []byte{})
			cr.stack = append(cr.stack, url)
		}
	}
}

// Crawler workers downloading, processing and saving sub URLs
func worker(cr *Crawler) {
	for url := range cr.input {
		fmt.Println(url)
		content := downloadWebPage(url)
		saveWebPage(cr, url, "saved", content)
		prepareChildWebPages(cr, cleanUrls(url, parseWebPage(url, string(content))))
		saveWebPage(cr, url, "done", content)
		cr.availableWorkers++
	}
}

// You can change the available workers in order to parralelize more processes
func main() {

	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Wrong arguments please specify a starting url and a destination folder\n")
		os.Exit(1)
	}

	cr := &Crawler{availableWorkers: 20, startUrl: os.Args[1], destinationPath: os.Args[2]}
	cr.init()

	for i := 0; i < cap(cr.input); i++ {
		go worker(cr)
	}

	for {
		if cr.availableWorkers > 0 && len(cr.stack) > 0 {
			cr.availableWorkers--
			cr.input <- cr.stack[0]
			cr.stack = cr.stack[1:]
		} else if cr.availableWorkers == 20 && len(cr.stack) == 0 {
			return
		}
	}
}
