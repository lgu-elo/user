include ./config/development.env

OUT_DIR=$(PWD)/pkg/pb
PROTO=docker run --rm -v "$(PWD):/defs" namely/protoc-all -f ./pkg/pb/user.proto -o ./ -l go

.SILENT: daemon mockgen protogen lint compose test

daemon:
	APP_ENV=development CompileDaemon \
		--build="go build -o ./bin/user ./cmd/user/main.go" \
		--command="./bin/user" \
		-graceful-kill \
		-log-prefix=false \
		-polling \
		-polling-interval=350

mockgen:
	go generate -v ./...

lint:
	golangci-lint run -c .golangcilint.yaml

compose:
	docker compose -f ./deployments/docker-compose.yml -p user up

protogen:
	@echo "Running generation protofiles..."
	@$(PROTO)

test:
	mkdir -p ./tmp
	go clean -testcache
	go test -v ./... -coverprofile=./tmp/coverage.out

build:
	GOARCH=amd64 \
	GOOS=linux \
	CGO_ENABLED=0 \
	APP_ENV=production \
	go build -o ./bin/user ./cmd/user/main.go

