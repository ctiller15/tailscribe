package api

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
	apiCfg := NewAPIConfig(NewEnvVars())
	apiCfg.HandleIndex(response, request)

	assert.Equal(t, response.Result().StatusCode, 200)
}

func TestGetAttributions(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/attributions", nil)
	response := httptest.NewRecorder()

	apiCfg := NewAPIConfig(NewEnvVars())
	apiCfg.HandleAttributions(response, request)

	assert.Equal(t, response.Result().StatusCode, 200)
}

func TestGetTerms(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/terms", nil)
	response := httptest.NewRecorder()

	apiCfg := NewAPIConfig(NewEnvVars())
	apiCfg.HandleTerms(response, request)

	assert.Equal(t, response.Result().StatusCode, 200)
}

func TestGetPrivacyPolicy(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/privacy", nil)
	response := httptest.NewRecorder()

	apiCfg := NewAPIConfig(NewEnvVars())
	apiCfg.HandlePrivacyPolicy(response, request)

	assert.Equal(t, response.Result().StatusCode, 200)
}
