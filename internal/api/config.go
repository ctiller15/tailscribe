package api

import "os"

type EnvVars struct {
	Addr         string
	ContactEmail string
}

func NewEnvVars() *EnvVars {
	addr := os.Getenv("PORT")
	contactEmail := os.Getenv("CONTACT_EMAIL")

	return &EnvVars{
		Addr:         addr,
		ContactEmail: contactEmail,
	}
}

type APIConfig struct {
	Env EnvVars
}

func NewAPIConfig(env *EnvVars) *APIConfig {
	return &APIConfig{
		Env: *env,
	}
}
