package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	if err := os.Chdir("../.."); err != nil {
		panic(err)
	}
}

func TestGetIndex(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	HandleIndex(response, request)

	assert.Equal(t, response.Result().StatusCode, 200)
}

func TestGetAttributions(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/attributions", nil)
	response := httptest.NewRecorder()

	HandleAttributions(response, request)

	assert.Equal(t, response.Result().StatusCode, 200)
}
