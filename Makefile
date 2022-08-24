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

get-lint:
	if [ ! -f ./bin/golangci-lint ]; then \
		wget -O - -q https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s latest; \
	fi;

lint: get-lint
	./bin/golangci-lint run ./... --timeout 5m0s

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
 
