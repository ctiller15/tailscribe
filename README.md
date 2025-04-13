# TailScribe

- A rewrite of [TailScribe](https://www.tailscribe.com/), a dog training notetaking app, written in Golang.

## Running the project

### Quickstart
```bash
cp .env.example .env
go run main.go
```

### Running the tests
```bash
go test ./...
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