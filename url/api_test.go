package url

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
)

var body *strings.Reader

var writer *httptest.ResponseRecorder
var request *http.Request
var withURLParam = "http://localhost:5555/hash?url=https://www.t-mobile.com/cell-phones"
var withOutURLParam = "http://localhost:5555/hash"
var withNoHostName = "http://localhost:5555/hash?url=www.t-mobile.com/cell-phones"

func setBodyWriterRequest(url string) {
	body = strings.NewReader("my request")
	writer = httptest.NewRecorder()
	request = httptest.NewRequest("GET", url, body)
}

func TestParamCheckNotFound(t *testing.T) {
	setBodyWriterRequest(withOutURLParam)
	value, found := paramCheck(writer, request, "user")
	assert.False(t, found, "expected found to equal false")
	assert.Empty(t, value, "expected value to be empty")
}

func TestParamCheckFound(t *testing.T) {
	setBodyWriterRequest(withURLParam)
	value, found := paramCheck(writer, request, "url")
	assert.True(t, found, "expected found to equal false")
	assert.NotEmpty(t, value, "expected value to be empty")
}

func TestHashHttp(t *testing.T) {
	os.Setenv("SALT", "ya")
	actual := Init()
	assert.NoError(t, actual, "exptected no error to be returned")
	setBodyWriterRequest(withURLParam)

	hashHTTP(writer, request)
	assert.NotEmpty(t, writer.Body, "expected body to not be emtpy")

	missingParam := missingParam{}
	json.Unmarshal(writer.Body.Bytes(), &missingParam)
	assert.Empty(t, missingParam.Msg, "expected missingParam message to be empty")

	parseURLErr := parseURLError{}
	json.Unmarshal(writer.Body.Bytes(), &parseURLErr)
	assert.Empty(t, parseURLErr.Msg, "expected parseURLErr message to be empty")

	assert.Equal(t, 200, writer.Code, "expected return code to equal 200")
}

func TestHashHttpNoURL(t *testing.T) {
	os.Setenv("SALT", "ya")
	actual := Init()
	assert.NoError(t, actual, "exptected no error to be returned")
	setBodyWriterRequest(withOutURLParam)
	hashHTTP(writer, request)
	assert.NotEmpty(t, writer.Body, "expected body to not be emtpy")

	missingParam := missingParam{}
	err := json.Unmarshal(writer.Body.Bytes(), &missingParam)
	assert.NotEmpty(t, missingParam.Msg, "expected missingParam message to be set")
	assert.Nil(t, err, "expected err to be nil")

	assert.Equal(t, 400, writer.Code, "expected return code to equal 400")
}

func TestHashHttpNoHostName(t *testing.T) {
	os.Setenv("SALT", "ya")
	actual := Init()
	assert.NoError(t, actual, "exptected no error to be returned")
	setBodyWriterRequest(withNoHostName)
	hashHTTP(writer, request)
	assert.NotEmpty(t, writer.Body, "expected body to not be emtpy")

	missingParam := missingParam{}
	err := json.Unmarshal(writer.Body.Bytes(), &missingParam)
	assert.NotEmpty(t, missingParam.Msg, "expected missingParam message to be set")
	assert.Nil(t, err, "expected err to be nil")

	assert.Equal(t, 400, writer.Code, "expected return code to equal 400")
}

func TestInitNoSalt(t *testing.T) {
	os.Unsetenv("SALT")
	actual := Init()
	assert.Error(t, actual, "expected error to be returned")
}

func TestInitWithSalt(t *testing.T) {
	os.Setenv("SALT", "ya")
	actual := Init()
	assert.NoError(t, actual, "exptected no error to be returned")
}
