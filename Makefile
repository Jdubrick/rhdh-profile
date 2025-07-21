BINARY_NAME=rhdh-cli

.PHONY: build
build:
	go build -o $(BINARY_NAME) .

.PHONY: deps
deps:
	go mod tidy
	go mod verify
.PHONY: deploy-operator-cli
deploy-operator-cli: build
	./$(BINARY_NAME) deploy operator --verbose

.PHONY: deploy-presets-cli
deploy-presets-cli: build
	./$(BINARY_NAME) deploy presets --verbose
