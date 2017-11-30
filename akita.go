// The MIT License (MIT)
//
// Copyright (c) 2017 Greivin López
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"net/url"

	"github.com/fatih/color"
	"github.com/jackdanger/collectlinks"
	"gopkg.in/fatih/set.v0"
)

const (
	versionNumber = "1.01"
)

var (
	website = ""
	version = false

	// Color objects for output printing
	green *color.Color
	red   *color.Color

	// Set to avoid infinite loops
	s *set.Set

	// Http client
	client http.Client
)

func fixUrl(href, base string) string {
	uri, err := url.Parse(href)
	if err != nil {
		return ""
	}
	baseUrl, err := url.Parse(base)
	if err != nil {
		return ""
	}
	uri = baseUrl.ResolveReference(uri)
	return uri.String()
}

func crawl(uri string) {
	printGreen := green.SprintFunc()
	printRed := red.SprintFunc()

	resp, err := client.Get(uri)
	if err != nil {
		fmt.Printf("Testing %s -> %s\n", uri, printRed("Fail: "+err.Error()))
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		fmt.Printf("Testing %s -> %s\n", uri, printGreen("OK"))
	} else {
		fmt.Printf("Testing %s -> %s\n", uri, printRed("Fail: "+resp.Status))
	}

	s.Add(uri)

	links := collectlinks.All(resp.Body)

	for _, link := range links {
		absolute := fixUrl(link, uri)
		if !s.Has(absolute) && uri != "" {
			u, err := url.Parse(absolute)
			if err != nil {
				panic(err)
			}
			if u.Host == "www.westernasset.com" {
				crawl(absolute)
			}
		}
	}
}

func init() {
	// Command line parameters
	flag.BoolVar(&version, "version", false, "Use this flag if you want to verify the version of the tool.")
	flag.StringVar(&website, "website", "", "Website url or host: http://www.google.com")

	// Color objects for output printing
	green = color.New(color.FgGreen)
	red = color.New(color.FgRed)

	// Creating http client
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client = http.Client{Transport: transport}

	// Initialize our thread safe Set
	s = set.New()
}

func main() {
	// Parse command line parameters
	flag.Parse()
	if version {
		fmt.Println("akita version " + versionNumber)
	} else {
		if website != "" {
			// Start crawling the website from the given root URL
			crawl(website)
		} else {
			fmt.Println("Please provide website: ./akita -website=\"http://www.google.com\" ")
		}
	}
}
