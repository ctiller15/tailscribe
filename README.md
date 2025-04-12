# TailScribe

- A rewrite of [TailScribe](https://www.tailscribe.com/), a dog training notetaking app, written in Golang.

## Running the project

### Quickstart
```bash
go run main.go
```

### Running the container
(Requires Docker)

```bash
# macos m1+
docker buildx build --platform linux/amd64 -t tailscribe .
docker run -p 8080:8080 tailscribe

# linux
docker build -t tailscribe .
docker run -p 8080:8080 tailscribe
```