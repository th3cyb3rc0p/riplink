package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/mschwager/riplink/src/parse"
	"github.com/mschwager/riplink/src/requests"
	"github.com/mschwager/riplink/src/rpurl"
)

func main() {
	var queryUrl string
	flag.StringVar(&queryUrl, "url", "https://google.com", "URL to query")

	var timeout int
	flag.IntVar(&timeout, "timeout", 5, "Timeout in seconds")

	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "Verbose output")

	flag.Parse()

	client := &http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}

	request, err := http.NewRequest("GET", queryUrl, nil)
	if err != nil {
		panic(err)
	}

	response, _, err := requests.SendRequest(client, request)
	if err != nil {
		panic(err)
	}

	node, err := parse.BytesToHtmlNode(response)
	if err != nil {
		panic(err)
	}

	anchors, err := parse.Anchors(node)
	if err != nil {
		panic(err)
	}

	hrefs, errs := parse.ValidHrefs(anchors)
	for _, err := range errs {
		fmt.Println("Invalid anchor:", err)
	}

	urls, errs := rpurl.AbsoluteHttpUrls(queryUrl, hrefs)
	for _, err := range errs {
		fmt.Println("Could not generate usable URL:", err)
	}

	var preparedRequests []*http.Request
	for _, url := range urls {
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println(err)
			continue
		}
		preparedRequests = append(preparedRequests, request)
	}

	results := make(chan *requests.Result)

	go requests.SendRequestsToChan(client, preparedRequests, results)

	for result := range results {
		if result.Err != nil {
			fmt.Println(result.Err)
			continue
		}

		if verbose || result.Code < 200 || result.Code > 299 {
			fmt.Println(result.Url, result.Code)
		}
	}
}
