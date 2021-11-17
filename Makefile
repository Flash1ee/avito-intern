LOG_DIR=./logs

stop-redis:
	systemctl stop redis
stop-postgres:
	systemctl stop postgresql

gen-mock:
	go generate ./...

run-coverage:
	go test -covermode=atomic -coverpkg=./internal/... -coverprofile=cover ./internal/...
	cat cover | fgrep -v "mocks" | fgrep -v "testing.go" | fgrep -v "docs"  | fgrep -v "configs" | fgrep -v "main.go" > cover2
	go tool cover -func=cover2

build: generate-api
	go build -o server.out -v ./cmd/server

generate-api:
	go get -u github.com/swaggo/swag/cmd/swag@v1.6.7
	go mod tidy
	swag init -g cmd/server/main.go -o docs

build-docker:
	docker build --no-cache --network host -f ./docker/Dockerfile . --tag balance-app

stop:
	docker-compose stop

rm-docker:
	docker rm -vf $$(docker ps -a -q) || true
run:
	mkdir -p $(LOG_DIR)
	docker-compose up --build --no-deps

cover-html:
	go test ./... -coverprofile cover.out
	go tool cover -html=cover.out -o cover.html

clean:
	rm -rf cover.html cover cover2 *.out
clean-logs:
	rm -rf ./logs/*.log

