package url

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"

	"github.com/patrickmn/go-cache"
)

//Input input to be hashed
type Input struct {
	URL string `json:"url"`
}

type missingParam struct {
	Msg string `json:"message"`
}

//Response hashed response
type Response struct {
	URL string `json:"url"`
}

type parseURLError struct {
	Msg string `json:"message"`
}

type notFoundInCache struct {
	Msg string `json:"message"`
}

//StartAPI run shortener API
func StartAPI() {
	http.HandleFunc("/hash", hashHTTP) // curl "localhost:5555/hash?url=https://www.t-mobile.com/cell-phones"
	http.HandleFunc("/open", openURL)  // curl "localhost:5555/open?url=http://localhost:5555/d2169e9"
	http.ListenAndServe(port, nil)
}

func hashHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)

	val, found := paramCheck(w, req, "url")
	if !found {
		return
	}

	urlParsed, err := url.Parse(val[0])
	if err != nil {
		e := parseURLError{Msg: err.Error()}
		enc.Encode(e)
		return
	}

	if len(urlParsed.Hostname()) == 0 {
		e := parseURLError{Msg: "unable to determine hostname:  the full url must be passed:  exmaple http://www.google.com"}
		enc.Encode(e)
		return
	}

	if urlParsed.RequestURI() == "/" {
		e := parseURLError{Msg: "unable to determine request URI:  the full url must be passed:  exmaple http://www.google.com/email"}
		enc.Encode(e)
		return
	}

	uri := urlParsed.RequestURI()
	uriTrim := uri[1:] //remove leading slash
	hashInput := Input{URL: uriTrim}
	result, _ := hash512(c, hashInput)
	hashedByte := result.([64]byte)
	hashedSubStr := fmt.Sprintf("%x", hashedByte[:4])
	// fmt.Printf("Hash512:\t%x\n\n", hashedByte[:4])

	resp := Response{
		URL: "http://localhost" + port + "/" + hashedSubStr,
	}
	c.Set(resp.URL, val[0], cache.DefaultExpiration) //add local_url[actual_url] to cache
	enc.Encode(resp)
}

func paramCheck(w http.ResponseWriter, req *http.Request, param string) ([]string, bool) {
	val, found := req.URL.Query()[param]
	if !found {
		m := missingParam{
			Msg: "Missing parameter:  " + param + " must be passed in querystring:  example localhost:5555/hash?url=http:www.googe.com",
		}
		enc := json.NewEncoder(w)
		enc.SetEscapeHTML(false)
		enc.Encode(m)
		return val, found
	}
	return val, found
}

/*
curl "localhost:5555/hash?url=https://www.t-mobile.com/cell-phones"
{"url":"http://localhost:5555/d2169e98"}

curl "localhost:5555/open?url=http://localhost:5555/d2169e98"
*/
func openURL(w http.ResponseWriter, req *http.Request) {
	localURL, found := paramCheck(w, req, "url")
	if !found {
		return
	}

	val, foundInCache := c.Get(localURL[0]) //look for localURL in cache
	if foundInCache {
		fmt.Println("found in cache")
	}

	if !foundInCache {
		n := notFoundInCache{
			Msg: "url not found in cache:  it must be submitted for shortening before it can be opened",
		}
		enc := json.NewEncoder(w)
		enc.SetEscapeHTML(false)
		enc.Encode(n)
		return
	}
	urlToOpen := val.(string)
	exec.Command("open", urlToOpen).Start() //open web url in browser
}
