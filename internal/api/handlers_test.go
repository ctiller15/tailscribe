package api

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/ctiller15/tailscribe/internal/database"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
)

var (
	TestEnvVars *EnvVars

	DbQueries *database.Queries
)

func teardown(ctx context.Context) {
	DbQueries.DeleteUsers(ctx)
}

func init() {
	ctx := context.Background()
	if err := os.Chdir("../.."); err != nil {
		panic(err)
	}

	if err := godotenv.Load(".env.test"); err != nil {
		log.Fatalf("error loading .env file: %v.\n", err)
	}

	TestEnvVars = NewEnvVars()

	dbUrl := TestEnvVars.Database.ConnectionString() + "?sslmode=disable"

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	DbQueries = database.New(db)

	defer teardown(ctx)
}

func TestGetIndex(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()
	apiCfg := NewAPIConfig(TestEnvVars, DbQueries)
	apiCfg.HandleIndex(response, request)

	assert.Equal(t, response.Result().StatusCode, 200)
}

func TestGetSignup(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/signup", nil)
	response := httptest.NewRecorder()
	apiCfg := NewAPIConfig(TestEnvVars, DbQueries)
	apiCfg.HandleSignupPage(response, request)

	assert.Equal(t, response.Result().StatusCode, 200)
}

func TestHandlePostSignup(t *testing.T) {
	t.Run("Happy path", func(t *testing.T) {
		formData := url.Values{
			"email":    {"invalidEmail@email.com"},
			"password": {"password123"},
		}

		request, _ := http.NewRequest(http.MethodPost, "/signup", strings.NewReader(formData.Encode()))
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()
		apiCfg := NewAPIConfig(TestEnvVars, DbQueries)
		apiCfg.HandlePostSignup(response, request)

		assert.Equal(t, 201, response.Result().StatusCode)
	})

	t.Run("Invalid email", func(t *testing.T) {
		formData := url.Values{
			"email":    {"invalidEmail"},
			"password": {"password123"},
		}

		request, _ := http.NewRequest(http.MethodPost, "/signup", strings.NewReader(formData.Encode()))
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()
		apiCfg := NewAPIConfig(TestEnvVars, DbQueries)
		apiCfg.HandlePostSignup(response, request)

		assert.Equal(t, response.Result().StatusCode, 400)
	})
}

func TestGetAttributions(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/attributions", nil)
	response := httptest.NewRecorder()

	apiCfg := NewAPIConfig(TestEnvVars, DbQueries)
	apiCfg.HandleAttributions(response, request)

	assert.Equal(t, response.Result().StatusCode, 200)
}

func TestGetTerms(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/terms", nil)
	response := httptest.NewRecorder()

	apiCfg := NewAPIConfig(TestEnvVars, DbQueries)
	apiCfg.HandleTerms(response, request)

	assert.Equal(t, response.Result().StatusCode, 200)
}

func TestGetPrivacyPolicy(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/privacy", nil)
	response := httptest.NewRecorder()

	apiCfg := NewAPIConfig(TestEnvVars, DbQueries)
	apiCfg.HandlePrivacyPolicy(response, request)

	assert.Equal(t, response.Result().StatusCode, 200)
}

func TestGetContactUs(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/contact", nil)
	response := httptest.NewRecorder()

	apiCfg := NewAPIConfig(TestEnvVars, DbQueries)
	apiCfg.HandleContactUs(response, request)

	assert.Equal(t, response.Result().StatusCode, 200)
}
