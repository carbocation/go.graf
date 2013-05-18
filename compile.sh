GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o forum.linux  *.go && ./nrsc-script forum.linux "static templates" -q
GOOS=darwin GOARCH=amd64 go build -o forum.osx  *.go && ./nrsc-script forum.osx "static templates" -q

