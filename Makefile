.PHONY: build

lint:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.44.2
	golangci-lint run ./...

vet:
	@go vet

test:
	@go test ./... -v

build:
	@go build .

docker:
	docker buildx build --platform "linux/amd64,linux/arm64,linux/arm/v5,linux/arm/v7" -t Am3o/co2_exporter:latest .

clean:
	@rm co2_exporter