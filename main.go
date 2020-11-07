package main

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

// Post information for sub reddit
type Post struct {
	ID        string    `json:"id"`
	Thumbnail Thumbnail `json:"media"`
}

// Thumbnail information.
type Thumbnail struct {
	Src string `json:"content"`
}

// Container for posts to live in.
type Container struct {
	Posts map[string]Post `json:"posts"`
}

func main() {
	url := "https://gateway.reddit.com/desktopapi/v1/subreddits/" + os.Args[1] + "?rtj=only&redditWebClient=web2x&app=web2x-client-production&allow_over18=1&include=prefsSubreddit&dist=7&layout=card&sort=hot&geo_filter=TR"
	var container Container
	getImages(&url, &container)
	downloadFiles(&container)
}

func getImages(url *string, cont *Container) {
	body := makeRequest(url)
	if err := json.Unmarshal([]byte(body), &cont); err != nil {
		log.Fatal(err)
	}
}

func makeRequest(url *string) []byte {
	var defaultTransport http.RoundTripper = &http.Transport{
		Proxy: nil,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          1,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   15 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	client := http.Client{Timeout: time.Second * 5, Transport: defaultTransport}

	req, reqErr := http.NewRequest("GET", *url, nil)
	if reqErr != nil {
		log.Fatalln(reqErr)
	}

	req.Header.Set("Host", "gateway.reddit.com")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:82.0) Gecko/20100101 Firefox/82.0")
	req.Header.Set("accept", "*/*")
	req.Header.Set("content-length", "0")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		panic(err.Error())
	}
	return body
}

func downloadFiles(cont *Container) {
	for k, v := range cont.Posts {
		if v.Thumbnail.Src != "" {
			downloadRequest(v.Thumbnail.Src, k)
		}
	}
}

func downloadRequest(URL, fileName string) error {
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}
	file, err := os.Create(fileName + ".jpg")
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	return nil
}
