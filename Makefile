.PHONY : help
help :
	@echo "Usage:"
	@echo
	@awk '/^###/ { $$1 = "   "; print }' ${MAKEFILE_LIST} | column -t -s ':'

### build: build for the current system
.PHONY : build
build :
	CGO_ENABLED=0 go build -v -trimpath

### install: build and install to `$GOPATH/bin` (`$HOME/go/bin` if the GOPATH environment variable is not set)
.PHONY : install
install :
	CGO_ENABLED=0 go install -v -trimpath

### cover: show coverage report in the default Web browser
.PHONY : cover
cover :
	go tool cover -html=cover.out

### test: run tests
.PHONY : test
test :
	go test ./src/...

### test-full: run tests with coverage
.PHONY : test-full
test-full :
	go test -v=1 -count=1 -coverprofile=cover.out ./src/...

### clean: remove binaries, coverage data, `bin` and `dist` folders
.PHONY : clean
clean :
	go clean
	rm -f *.out
	rm -rf bin/
	rm -rf dist/
