GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o forum.linux  main.go
GOOS=darwin GOARCH=amd64 go build -o forum.osx  main.go

