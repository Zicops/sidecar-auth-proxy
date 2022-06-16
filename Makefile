.PHONY: default mac linux windows docker clean
BINARY_NAME = sidecar-auth-proxy
CHART_NAME = sidecar-auth-proxy
default: linux

all: linux test

docker: linux
	@echo "building docker image" ;\
		docker build -t "$(BINARY_NAME):localdeploy" .
clean:
	-rm $(BINARY_NAME)
mac:
	@echo "building $(BINARY_NAME) (mac)" ;\
        go build -o $(BINARY_NAME) ./cmd/sidecar-auth-proxy
linux:
	## CGO_ENABLED=0 go build -a -installsuffix cgo is not needed for Go 1.10 or later
	## https://github.com/golang/go/issues/9344#issuecomment-69944514
	@echo "building $(BINARY_NAME) (linux)" ;\
        GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME) ./cmd/sidecar-auth-proxy
windows:
	@echo "building $(BINARY_NAME) (windows)" ;\
      		env GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME) ./cmd/sidecar-auth-proxy
test:
	@CC=gcc go test ./...

test-race:
	@CC=gcc go test -race -short ./...

validate-circleci:
	@circleci config validate -c .circleci/config.yml

test:
	@CC=gcc go test ./...

test-race:
	@CC=gcc go test -race -short ./...

deps-lint:
	@GO111MODULE=off go get golang.org/x/lint
	@GO111MODULE=off go get golang.org/x/lint/golint
	@GO111MODULE=off go get golang.org/x/tools/cmd/goimports
	@GO111MODULE=off go get golang.org/x/tools/go/analysis/passes/nilness/cmd/nilness

deps-verify:
	@go mod tidy
	@go mod verify

lint:
	@go vet ./...
	@go vet -vettool=$(which nilness) ./...
	@go fix ./...
	@golint ./...
	# @! goimports -l . | grep -vF 'No Exceptions'
testandcover:
	@echo "Getting go-iunit-report and gocov and gocov-xml for testing"
	@GO111MODULE=off go get github.com/jstemmer/go-junit-report
	@GO111MODULE=off go get github.com/axw/gocov/gocov
	@GO111MODULE=off go get gopkg.in/matm/v1/gocov-html
	@GO111MODULE=off go get github.com/AlekSi/gocov-xml
	@echo "clean old test files"
	@rm -f Tests-*
	@rm -rf coverage/
	@mkdir coverage

	echo "Start go test"
	@GO111MODULE=off go clean -testcache
	@go test -v -coverprofile=coverage/coverage.out -covermode count ./... | go-junit-report > Tests-$(go env GOOS)-report.xml
	@exitCodeTests=${PIPESTATUS[0]}
	@echo "Exit code tests: $exitCodeTests"

	@echo "Start calculating code coverage"
	@gocov convert coverage/coverage.out > coverage/coverage.json
	@gocov-xml < coverage/coverage.json > coverage/coverage.xml
	@gocov-html < coverage/coverage.json > coverage/index.html
	@exit $exitCodeTests

