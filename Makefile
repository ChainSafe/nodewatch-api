PROJECTNAME=$(shell basename "$(PWD)")
GOLANGCI := $(GOPATH)/bin/golangci-lint

.PHONY: help lint test run
all: help
help: Makefile
	@echo
	@echo " Choose a make command to run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

.PHONY: get-lint
get-lint:
	if [ ! -f ./bin/golangci-lint ]; then \
		curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.54.2; \
	fi;

.PHONY: lint
lint: get-lint
	@echo "  >  \033[32mRunning lint...\033[0m "
	./bin/golangci-lint run --config=./.golangci.yml

## license: Adds license header to missing files.
license:
	@echo "  >  \033[32mAdding license headers...\033[0m "
	GO111MODULE=off go get -u github.com/google/addlicense
	addlicense -c "ChainSafe Systems" -f ./copyright.txt -y 2021 .

## license-check: Checks for missing license headers
license-check:
	@echo "  >  \033[32mChecking for license headers...\033[0m "
	GO111MODULE=off go get -u github.com/google/addlicense
	addlicense -check -c "ChainSafe Systems" -f ./copyright.txt -y 2021 .

test:
	go test ./...

build:
	go build -o ./bin/crawler cmd/main.go

run:
	@echo "  >  \033[32mUsing Docker Container for development...\033[0m "
	@echo "  >  \033[32mRemoving old User Service stuff...\033[0m "
	docker-compose -f docker-compose.yaml down -v
	@echo "  >  \033[32mStarting Crawler Service w/o build...\033[0m "
	docker-compose -f docker-compose.yaml up --build $$scale
 
