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

test:
	go test ./...

run:
	@echo "  >  \033[32mUsing Docker Container for development...\033[0m "
	@echo "  >  \033[32mRemoving old User Service stuff...\033[0m "
	docker-compose -f docker-compose.yaml down -v
	@echo "  >  \033[32mStarting Crawler Service w/o build...\033[0m "
	docker-compose -f docker-compose.yaml up
 