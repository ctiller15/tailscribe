source .env
dbconnstr="postgresql://${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?user=${POSTGRES_USER}&password=${POSTGRES_PASSWORD}"

cd sql/schema 
goose postgres $dbconnstr down