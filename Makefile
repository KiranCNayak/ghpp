.PHONY: build clean

APP_NAME := ghpp

install:
	@echo "Installing dependencies..."
	@go install
	@echo "Dependencies installed!"

build: clean
	@echo "Building the project..."
	@go build -o $(APP_NAME) .
	@echo "Build completed!"

clean:
	@rm -rf $(APP_NAME)

