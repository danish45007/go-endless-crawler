package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/steelx/extractlinks"
)

func checkErrors(err error) {
	if(err != nil) {
		fmt.Println("Error",err)
		os.Exit(1)
	}
}


var (
	// custom client to bypass ssl certificate
	config = &tls.Config{
		InsecureSkipVerify: true,
	}
	transport = &http.Transport{
		TLSClientConfig: config,
	}
	netClient = &http.Client{
		Transport: transport,
	}
	queue = make(chan string)
	hasVisited = make(map[string]bool)
)

func validUrl(href,baseUrl string) string {
	uri,err := url.Parse(href)
	if(err != nil) {
		return ""
	}

	base,err := url.Parse(baseUrl)
	if(err != nil) {
		return ""
	}
	// base.Host + uri.Path
	fixedURI := base.ResolveReference(uri)
	return fixedURI.String()
}

func isSameDomain(href,baseUrl string) bool {
	uri,err := url.Parse(href)
	if(err != nil) {
		return false
	}

	base,err := url.Parse(baseUrl)
	if(err != nil) {
		return false
	}
	
	if uri.Host != base.Host {
		return false
	}

	return true
}


func urlPraser(href string) {
	fmt.Printf("Currently Crawling ---> %v \n",href)
	hasVisited[href] = true
	response, err := netClient.Get(href)
	checkErrors(err)
	defer response.Body.Close()
	links, err := extractlinks.All(response.Body)
	checkErrors(err)
	for _, link := range links {
		absLink := validUrl(link.Href,href)
		go func() {
			queue <- absLink
		}()
	}
}

func Crawler() {
	args := os.Args[1:]
	if(len(args) == 0) {
		fmt.Println("Please enter a url")
		os.Exit(1)
	} 
	_, err := url.ParseRequestURI(args[0])
	if(err != nil) {
		fmt.Println("Please enter a valid url")
	}
	baseUrl := args[0]
	// go-routine concurrent exec.
	go func (){
		// send base-url into queue
		queue <- baseUrl
	}()

	for href := range queue {
		// not visited url's
		if(!hasVisited[href]) && isSameDomain(href,baseUrl) {
			urlPraser(href)
		}
	}
}


func main() {
	Crawler()
}