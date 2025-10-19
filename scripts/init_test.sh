# Intended for local development use.
# Initializes the local test environment

source .env.test

# Name and port defined in docker-compose
dbconnstr="postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}"

cd sql/schema 
goose postgres $dbconnstr up