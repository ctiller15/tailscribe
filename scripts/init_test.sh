# Intended for local development use.
# Initializes the local test environment

source .env
# Name and port defined in docker-compose
test_postgres_port=6432
test_postgres_db="animal_training_journal_test"
dbconnstr="postgresql://${POSTGRES_HOST}:${test_postgres_port}/${test_postgres_database}?user=${POSTGRES_USER}&password=${POSTGRES_PASSWORD}"

cd sql/schema 
goose postgres $dbconnstr up