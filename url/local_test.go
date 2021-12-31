package url

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
)

func setup() {
	os.Setenv("SALT", "helloworld")
	if c == nil {
		c = cache.New(10*time.Minute, 10*time.Minute)
	}
}

func TestSha512(t *testing.T) {
	setup()

	txt := "Hello World!"
	expectedHash := "fc3403cd"

	i := Input{URL: txt}

	val512, _ := hash512(c, i)
	val := val512.([64]byte)
	actualHash := fmt.Sprintf("%x", val[:4])
	assert.Equal(t, expectedHash, actualHash)
}

func TestSha512Negative(t *testing.T) {
	setup()

	txt := "http://www.google.com"
	expectedHash := "861844d6"

	i := Input{URL: txt}

	val512, _ := hash512(c, i)
	val := val512.([64]byte)
	actualHash := fmt.Sprintf("%x", val[:4])
	assert.NotEqual(t, expectedHash, actualHash)
}

func TestSha512One(t *testing.T) {
	setup()

	txt := "http://www.google.com"
	expectedHash := "9925924f"

	i := Input{URL: txt}

	val512, _ := hash512(c, i)
	val := val512.([64]byte)
	actualHash := fmt.Sprintf("%x", val[:4])
	assert.Equal(t, expectedHash, actualHash)
}

func TestSha512Cache(t *testing.T) {
	setup()

	txt := "Hello Cache"
	expectedHash := "730e161f"

	i := Input{URL: txt}

	val512, found := hash512(c, i)
	val := val512.([64]byte)
	actualHash := fmt.Sprintf("%x", val[:4])
	assert.Equal(t, expectedHash, actualHash)
	assert.False(t, found, "expected found to be false")

	val512, found = hash512(c, i)
	val = val512.([64]byte)
	actualHash = fmt.Sprintf("%x", val[:4])
	assert.Equal(t, expectedHash, actualHash)
	assert.True(t, found, "expected found to be true")
}
