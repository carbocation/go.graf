GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main.linux  main.go
GOOS=darwin GOARCH=amd64 go build -o main.osx  main.go

