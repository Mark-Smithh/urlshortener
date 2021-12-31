package url

import (
	"crypto/sha512"
	"fmt"
	"os"
	"time"

	"github.com/patrickmn/go-cache"
)

var c *cache.Cache
var port = ":5555"

//Init initialize needed objects
func Init() {
	found, err := checkSalt()
	if !found && err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	c = cache.New(10*time.Minute, 10*time.Minute)
}

//Hash512256 hash using SHA512/256
func hash512256(txt string) [32]byte {
	sha512 := sha512.Sum512_256([]byte(txt))
	return sha512
}

//Hash512 hash using SHA512
func hash512(c *cache.Cache, input Input) (interface{}, bool) {
	val, foundInCache := c.Get(input.URL)
	if foundInCache {
		// fmt.Println("found in cache")
		return val, foundInCache
	}

	foundSalt, err := checkSalt()
	if !foundSalt && err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	salt, _ := os.LookupEnv("SALT")
	salted := input.URL + salt
	// fmt.Printf("salted: %s\n", salted)
	sha512 := sha512.Sum512([]byte(salted))
	c.Set(input.URL, sha512, cache.DefaultExpiration) //add url[sha512] to cache
	return sha512, foundInCache
}

func checkSalt() (bool, error) {
	var initErr error
	_, found := os.LookupEnv("SALT")
	if !found {
		initErr = fmt.Errorf("required environment variable SALT not found")
		return found, initErr
	}
	return found, nil
}
