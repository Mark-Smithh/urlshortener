package main

import (
	"os"
	"strings"
	"urlshortener/url"
)

func main() {
	url.Init()
	if len(os.Args) >= 2 && os.Args[1] == strings.ToLower("local") {
		url.Local()
	}
	url.StartAPI()
}
