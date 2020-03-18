run:
	docker-compose up -d

run-build:
	docker-compose up -d --build

down:
	docker-compose down

unit-test:
	go test ./... -count=1

unit-test-cover:
	go list ./... | xargs go test -v -cover -coverprofile=c.out
	go tool cover -html=c.out -o coverage.html