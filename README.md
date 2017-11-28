# akita
Akita is a command line tool to crawl a website looking for broken links.

The tool is named after the [Akita dog breed](https://en.wikipedia.org/wiki/Akita_(dog)) which are very protective of their home & owners.

![Akira dog](https://i.pinimg.com/564x/00/51/3e/00513e9bad23301d989740c8ca266a91.jpg "Akira dog").

## Installation

The akita tool is written using [Go](https://golang.org/).  So the first step is to download and install Go and set your development environment.

Akita uses packages outside of the standard library, each of those packages need to be imported on the go environment before compiling the tool:

```
go get github.com/fatih/color
go get gopkg.in/fatih/set.v0
go get github.com/jackdanger/collectlinks
```
Download the code and put it on a folder named 'akita' inside the $GOPATH/src folder of your Go environment.

Compile the tool:
```
go build akita.go
```
## Run the tool

To run the tool call it as:
```
./akita -website="http://www.yoursite.com"
```
