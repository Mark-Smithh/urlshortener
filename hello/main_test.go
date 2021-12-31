package hello

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpper(t *testing.T) {
	result := Upper("hello world")
	expected := "HELLO WORLD"
	assert.Equal(t, expected, result)
}
