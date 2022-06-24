build-image:
	docker build --file=Dockerfile --tag=gurygury/test-hello:latest .

start-container:
	docker run --env-file .env --name=test-hello -p 8080:8080 gurygury/test-hello:latest

test:
	go test ./... -coverprofile cover.out.tmp

lint:
	golangci-lint run --config=.golangci.yml