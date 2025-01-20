APP_NAME=mstream-api
BUILD_DIR="./build/$(APP_NAME)"

install:
	go get -u ./... && go mod tidy

run:
	nodemon --exec "go run" .

build:
	mkdir -p ./build && CGO_ENABLED=0 GOOS=linux go build -o ${BUILD_DIR}

test:
	go test -cover -v ./...

generate:
	go generate ./...