BINARY_NAME=Projectapirest
GO=go

all: build

run:
	@echo "Running the project..." #работает
	$(GO) run ./cmd

test:
	@echo "Running tests..."  #работает
	$(GO) test ./...

migrate-up:
	@echo "Applying migrations..." #работает
	goose up