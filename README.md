# TailScribe

- A rewrite of [TailScribe](https://www.tailscribe.com/), a dog training notetaking app, written in Golang.

## Running the project

### Installations
- [sqlc](https://docs.sqlc.dev/en/stable/overview/install.html)
- [docker desktop](https://docs.docker.com/desktop/)
- [Golang 1.24+](https://go.dev/doc/install)

### Quickstart
```bash
# copy env files
cp .env.example .env

# start the database
docker-compose up -d

# start server
go run main.go
```

### Running the tests
```bash
# Start up docker database container
docker-compose up animal-training-journal-test-db

# Run local migrations - get container up to date
./scripts/init_test.sh

# Run tests
go test ./...

# Clearing the database
rm -rf ./postgres/test_data
```

### connecting to the db manually
```bash
# Local
psql -U postgres -h localhost

# Test db
psql postgresql://postgres:postgres@localhost:6432/animal_training_journal_test
```

### running goose migrations
```bash
# move up migrations
./scripts/goose_up.sh

# move down migrations
./scripts/goose_down.sh
```

### Running the container
(Requires Docker)

```bash
# macos m1+
docker buildx build --platform linux/amd64 -t tailscribe .

# linux
docker build -t tailscribe .

# run
docker run --env-file=.env -p 8080:8080 tailscribe
```