package main

import (
	"fmt"
	"os"
	"strings"
	"urlshortener/url"
)

func main() {
	err := url.Init()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if len(os.Args) >= 2 && os.Args[1] == strings.ToLower("help") {
		fmt.Println("runs on port 5555")
		fmt.Println("to shorten a URL run: curl \"localhost:5555/hash?url=https://www.t-mobile.com/cell-phones\" returns [SHORTED_URL]")
		fmt.Println("to open a URL run: curl \"localhost:5555/open?url=[SHORTED_URL]\"")
		return
	}

	if len(os.Args) >= 2 && os.Args[1] == strings.ToLower("local") {
		url.Local()
	}
	url.StartAPI()
}
