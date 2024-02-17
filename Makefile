.PHONY: build run clean

build:
	@echo "Building Go Lambda function"
	@docker run --rm -v $(PWD):/app -w /app -e GOOS=linux -e GOARCH=amd64 golang:latest go build
deploy:
	@scp mail_gateway bridgenote@103.150.190.170:/home/bridgenote/career-bridge/mail_gateway                                   