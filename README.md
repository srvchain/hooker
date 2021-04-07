# Hooker

## Build

go build \
  -ldflags "-X main.version=${VERSION}"
  -o hooker \
  cmd/server/server.go

### cross compile linux

ENV="CGO_ENABLED=0 GOOS=linux GOARCH=amd64" \
  go build \
  -ldflags "-X main.version=${VERSION}"
  -o hooker \
  cmd/server/server.go

## RUN

./hooker

### docker run
	docker run -p 80:80 \
	  -v ${pwd}:/app -w /app \
	  --restart unless-stopped --name hooker -d debian ./hooker