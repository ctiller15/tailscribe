# Intended for local development use.
# Initializes the local test environment

if [ -f .env.test ]; then
    source .env.test
fi

# Name and port defined in docker-compose
dbconnstr="postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}"

cd sql/schema 
goose postgres $dbconnstr up