default: check-license test lint

test:
	go test ./... -coverprofile=cp.out

lint:
	golangci-lint run

build:
	go build

check-license:
	./scripts/check_license.sh
