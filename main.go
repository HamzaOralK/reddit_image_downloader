package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

func main() {
	fmt.Println("https://www.reddit.com/r/" + os.Args[1] + "/top/?t=day")
	url := "https://www.reddit.com/r/" + os.Args[1] + "/top/?t=day"

	c := make(chan string)

	go func(v string) {
		getImages(url, c)
		defer close(c)
	}(url)

	for v := range c {
		fmt.Println(v)
	}
}

func getImages(url string, channel chan string) {
	c := colly.NewCollector()

	// On every a element which has href attribute call callback
	c.OnHTML("img", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("class"), "_1XWObl-3b9tPy64oaG6fax") {
			imageLink := e.Attr("src")
			// fmt.Println("Image found. with src=", imageLink)
			channel <- imageLink
		}
	})
	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	// Start scraping on https://hackerspaces.org
	c.Visit(url)
}
