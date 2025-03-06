all:
	go run main.go

fetch:
	# mkdir -p /home/me/projects && cd /home/me/projects
	# git clone https://github.com/paepckehh/keys	
	GITSTORE="/home/me/projects/keys" go run main.go fetch 

check: 
	gofmt -w -s .
	go vet .
	golangci-lint run
	staticcheck
