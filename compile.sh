GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOMAXPROCS=8 go build -o forum.linux  *.go && ./nrsc-script forum.linux "static templates" -q
GOOS=darwin GOARCH=amd64 GOMAXPROCS=4 go build -o forum.osx  *.go && ./nrsc-script forum.osx "static templates" -q

