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

	err := godotenv.Load(".env.test")

	if err != nil {
		log.Printf("error loading .env file: %v.\n May experience degraded behavior during tests.\n", err)
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

func signUserUp(email, password string) []*http.Cookie {
	formData := url.Values{
		"email":    {email},
		"password": {password},
	}

	request, _ := http.NewRequest(http.MethodPost, "/signup", strings.NewReader(formData.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()
	apiCfg := NewAPIConfig(TestEnvVars, DbQueries)
	apiCfg.HandlePostSignup(response, request)

	result := response.Result()

	cookies := result.Cookies()
	return cookies
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
			"email":    {"testEmail@email.com"},
			"password": {"password123"},
		}

		request, _ := http.NewRequest(http.MethodPost, "/signup", strings.NewReader(formData.Encode()))
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()
		apiCfg := NewAPIConfig(TestEnvVars, DbQueries)
		apiCfg.HandlePostSignup(response, request)

		result := response.Result()
		assert.Equal(t, 302, result.StatusCode)

		assert.Equal(t, result.Header.Get("Location"), "/add_new_pet")

		cookies := result.Cookies()
		assert.NotNil(t, cookies[0])
		assert.Equal(t, "token", cookies[0].Name)
		assert.NotNil(t, cookies[1])
		assert.Equal(t, "refresh_token", cookies[1].Name)
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

func TestHandleLogout(t *testing.T) {
	t.Run("Logs out if currently logged in", func(t *testing.T) {
		formData := url.Values{
			"email":    {"testEmail1@email.com"},
			"password": {"password123"},
		}

		request, _ := http.NewRequest(http.MethodPost, "/signup", strings.NewReader(formData.Encode()))
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()
		apiCfg := NewAPIConfig(TestEnvVars, DbQueries)
		apiCfg.HandlePostSignup(response, request)

		result := response.Result()

		cookies := result.Cookies()

		auth_cookie := *cookies[0]

		logoutRequest, _ := http.NewRequest(http.MethodPost, "/logout", strings.NewReader(""))
		logoutRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		logoutRequest.AddCookie(&auth_cookie)
		logoutResponse := httptest.NewRecorder()
		apiCfg.HandlePostLogout(logoutResponse, logoutRequest)

		logoutResult := logoutResponse.Result()

		assert.Equal(t, 302, logoutResult.StatusCode)
		logoutCookies := logoutResult.Cookies()
		assert.NotNil(t, logoutCookies[0])
		assert.Equal(t, "token", logoutCookies[0].Name)
		assert.Equal(t, "", logoutCookies[0].Value)

		assert.Equal(t, "refresh_token", logoutCookies[1].Name)
		assert.Equal(t, "", logoutCookies[1].Value)
	})

	t.Run("Successfully logs out even if not logged in", func(t *testing.T) {
		apiCfg := NewAPIConfig(TestEnvVars, DbQueries)

		logoutRequest, _ := http.NewRequest(http.MethodPost, "/logout", strings.NewReader(""))
		logoutRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		logoutResponse := httptest.NewRecorder()
		apiCfg.HandlePostLogout(logoutResponse, logoutRequest)

		logoutResult := logoutResponse.Result()

		assert.Equal(t, 302, logoutResult.StatusCode)
		logoutCookies := logoutResult.Cookies()
		assert.NotNil(t, logoutCookies[0])
		assert.Equal(t, "token", logoutCookies[0].Name)
		assert.Equal(t, "", logoutCookies[0].Value)

		assert.Equal(t, "refresh_token", logoutCookies[1].Name)
		assert.Equal(t, "", logoutCookies[1].Value)
	})
}

func TestHandlePostLogin(t *testing.T) {
	t.Run("Successfully logs a user in", func(t *testing.T) {
		formData := url.Values{
			"email":    {"testEmail2@email.com"},
			"password": {"password123"},
		}

		request, _ := http.NewRequest(http.MethodPost, "/signup", strings.NewReader(formData.Encode()))
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()
		apiCfg := NewAPIConfig(TestEnvVars, DbQueries)
		apiCfg.HandlePostSignup(response, request)

		result := response.Result()
		assert.Equal(t, 302, result.StatusCode)

		loginFormData := url.Values{
			"email":    {"testEmail2@email.com"},
			"password": {"password123"},
		}

		loginRequest, _ := http.NewRequest(http.MethodPost, "/login", strings.NewReader(loginFormData.Encode()))
		loginRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		loginResponse := httptest.NewRecorder()
		loginApiCfg := NewAPIConfig(TestEnvVars, DbQueries)
		loginApiCfg.HandlePostLogin(loginResponse, loginRequest)

		loginResult := loginResponse.Result()
		assert.Equal(t, 302, loginResult.StatusCode)

		assert.Equal(t, loginResult.Header.Get("Location"), "/dashboard")

		cookies := loginResult.Cookies()
		assert.NotNil(t, cookies[0])
		assert.Equal(t, "token", cookies[0].Name)
		assert.NotNil(t, cookies[1])
		assert.Equal(t, "refresh_token", cookies[1].Name)
	})

	t.Run("Fails to log a user in if invalid username/password", func(t *testing.T) {
		loginFormData := url.Values{
			"email":    {"testEmail3@email.com"},
			"password": {"password123"},
		}

		loginRequest, _ := http.NewRequest(http.MethodPost, "/login", strings.NewReader(loginFormData.Encode()))
		loginRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		loginResponse := httptest.NewRecorder()
		loginApiCfg := NewAPIConfig(TestEnvVars, DbQueries)
		loginApiCfg.HandlePostLogin(loginResponse, loginRequest)

		loginResult := loginResponse.Result()
		assert.Equal(t, 401, loginResult.StatusCode)

		cookies := loginResult.Cookies()
		assert.Len(t, cookies, 0)
	})
}

// Authorized routes
func TestGetAddNewPethandler(t *testing.T) {
	t.Run("Fails to find add new pet page when unauthorized", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/dashboard/add_new_pet", nil)
		response := httptest.NewRecorder()

		apiCfg := NewAPIConfig(TestEnvVars, DbQueries)
		apiCfg.CheckAuthMiddleware(
			apiCfg.HandleAddNewPet,
		)(response, request)

		result := response.Result()

		assert.Equal(t, 401, result.StatusCode)
		assert.Equal(t, result.Header.Get("Location"), "/login")
	})

	t.Run("Succeeds at finding add new pet page when authorized", func(t *testing.T) {
		cookies := signUserUp("testEmail4@test.com", "password123")

		auth_cookie := *cookies[0]

		request, _ := http.NewRequest(http.MethodGet, "/dashboard/add_new_pet", nil)
		request.AddCookie(&auth_cookie)

		response := httptest.NewRecorder()
		apiCfg := NewAPIConfig(TestEnvVars, DbQueries)
		apiCfg.CheckAuthMiddleware(
			apiCfg.HandleAddNewPet,
		)(response, request)

		result := response.Result()
		assert.Equal(t, 200, result.StatusCode)
	})
}

// Unauthorized routes
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
