.SILENT:

lint:
	golangci-lint run ./...

build:
	docker build -t fetcher .

run:
	docker run --name=fetcher --rm -t -i fetcher
