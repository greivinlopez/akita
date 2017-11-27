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
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/smallnest/goreq"
	"time"
)

const (
	versionNumber    = "1.0"
)

var (
	website    = ""
	version   = false
)

func crawl(url string) {
    // Creating http client
	var request *goreq.GoReq
	request = goreq.New().Timeout(7 * time.Second)

    green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

    resp, _, _ := request.Get(url).EndBytes()

    if resp.StatusCode == 200 {
        fmt.Printf("Testing %s -> %s\n", url, green("OK"))
    } else {
        fmt.Printf("Testing %s -> %s\n", url, red("Fail: "+resp.Status))
    }
}

func init() {
	flag.BoolVar(&version, "version", false, "Use this flag if you want to verify the version of the tool.")
	flag.StringVar(&website, "website", "http://www.westernasset.com/us/en/", "Website url or host: http://www.westernasset.com/us/en/.")
}

func main() {
	// Parse command line parameters
	flag.Parse()
	if version {
		fmt.Println("akita version " + versionNumber)
	} else {
		crawl(website)
	}
}
