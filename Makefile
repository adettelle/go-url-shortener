test:
	go test ./...

lint: 
	golangci-lint run

vet:
	staticcheck ./...

check: lint vet test

testcov: test
	go test -v -coverpkg=./... -coverprofile=profile.cov ./... && go tool cover -func profile.cov && go tool cover -html profile.cov

testcov2: test
	go test -v -coverpkg=./... -coverprofile=profile.cov ./... && \
	cat profile.cov | grep -v "mock_.*.go" > cover.out && \
	go tool cover -func cover.out && go tool cover -html cover.out -o /mnt/c/Users/llitv/Documents/cover.html && \
	wslview file:///C:/Users/llitv/Documents/cover.html
	
run-server:
	go run ./cmd/shortener/