package api

import (
	"fmt"
	"os"

	"github.com/ctiller15/tailscribe/internal/database"
)

type DatabaseEnv struct {
	Name     string
	User     string
	Password string
	Host     string
	Port     string
}

type EnvVars struct {
	Addr         string
	ContactEmail string
	Database     DatabaseEnv
}

func NewEnvVars() *EnvVars {
	addr := os.Getenv("PORT")
	contactEmail := os.Getenv("CONTACT_EMAIL")
	dbName := os.Getenv("POSTGRES_DATABASE")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")

	return &EnvVars{
		Addr:         addr,
		ContactEmail: contactEmail,
		Database: DatabaseEnv{
			Name:     dbName,
			User:     dbUser,
			Password: dbPassword,
			Host:     dbHost,
			Port:     dbPort,
		},
	}
}

type APIConfig struct {
	Env EnvVars
	Db  database.Queries
}

func NewAPIConfig(env *EnvVars, db *database.Queries) *APIConfig {
	return &APIConfig{
		Env: *env,
		Db:  *db,
	}
}

// Creates the connection string for the database instance
func (d *DatabaseEnv) ConnectionString() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.Name,
	)
}
