.SILENT:

lint:
	golangci-lint run ./...

# Only run if you have five minutes to spare.
build:
	docker build -t delivery-service .

build-dev:
	docker build -t delivery-service -f Dockerfile.dev .

run:
	docker run --rm -t -i delivery-service
