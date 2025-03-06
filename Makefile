PROJECT=$(shell basename $(CURDIR))

all:
	make -C cmd/$(PROJECT) all

fetch:
	make -C cmd/$(PROJECT) fetch

deps: 
	rm go.mod 
	go mod init paepcke.de/$(PROJECT)
	go mod tidy -v	

check: 
	gofmt -w -s .
	go vet .
	staticcheck
	golangci-lint run
	make -C cmd/$(PROJECT) check
