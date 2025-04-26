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
go test ./...
```

### connecting to the db manually
```bash
psql -U postgres -h localhost
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