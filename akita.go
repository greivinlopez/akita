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
	"encoding/csv"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/jackdanger/collectlinks"
	"gopkg.in/fatih/set.v0"
)

const (
	versionNumber = "1.05"
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

	info *log.Logger

	// Slice of broken link reports
	errors []BrokenLink
)

// BrokenLink represents the report of an error found when testing GET
// a given URL. The tool will crawl the website and each time it founds
// a problem it will create a new BrokenLink and stores it into a collection
// of reports that will eventually be serialized into a report file.
type BrokenLink struct {
	Origin       string
	Target       string
	ErrorMessage string
}

// GetHeaders returns an array of strings containing the header titles
// intended for a CSV serialization of a collection of BrokenLink structures.
func getHeaders() []string {
	return []string{"Origin", "Target", "Error"}
}

// GetValues returns an array of strings containing the values of a given
// BrokenLink report. The result is intended to be serialized on a CSV file.
func (broken *BrokenLink) getValues() []string {
	return []string{broken.Origin, broken.Target, broken.ErrorMessage}
}

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

func reportError(origin string, uri string, message string) {
	broken := BrokenLink{Origin: origin, Target: uri, ErrorMessage: message}
	errors = append(errors, broken)
}

func crawl(origin string, uri string) {
	printGreen := green.SprintFunc()
	printRed := red.SprintFunc()

	resp, err := client.Get(uri)
	if err != nil {
		info.Printf("Testing %s -> %s", uri, printRed("Fail: "+err.Error()))
		reportError(origin, uri, err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		info.Printf("Testing %s -> %s", uri, printGreen("OK"))
	} else {
		info.Printf("Testing %s -> %s", uri, printRed("Fail: "+resp.Status))
		reportError(origin, uri, "Fail: "+resp.Status)
	}

	s.Add(uri)

	links := collectlinks.All(resp.Body)

	for _, link := range links {
		absolute := fixUrl(link, uri)
		if !s.Has(absolute) && uri != "" {
			u, err := url.Parse(absolute)
			if err != nil {
				info.Printf(err.Error())
			}
			if u.Host == "www.westernasset.com" {
				crawl(uri, absolute)
			}
		}
	}
}

func createReport() error {
	errorsCount := len(errors)
	records := make([][]string, errorsCount)

	for i := 0; i < errorsCount; i++ {
		records[i] = errors[i].getValues()
	}

	// write the file
	f, err := os.Create("result.csv")
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)

	// Write headers
	headers := getHeaders()
	if err = w.Write(headers); err != nil {
		return err
	}

	w.WriteAll(records) // calls Flush internally

	if err := w.Error(); err != nil {
		return err
	}
	return nil
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

	info = log.New(os.Stdout, "", 0)
}

func main() {
	// Store start time
	start := time.Now()

	// Parse command line parameters
	flag.Parse()
	if version {
		info.Println("akita version " + versionNumber)
	} else {
		if website != "" {
			// Start crawling the website from the given root URL
			crawl(website, website)
		} else {
			info.Println("Please provide website: ./akita -website=\"http://www.google.com\" ")
		}
	}

	// Output total time elapsed
	info.Println("Task completed")
	info.Printf("%.2fs elapsed\n", time.Since(start).Seconds())

	if len(errors) > 0 {
		createReport()
	}
}
