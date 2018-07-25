VERSION=$(shell (git describe --abbrev=0 --tags || echo '0.1.0') 2>/dev/null)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
HEAD=$(shell (git rev-list --abbrev-commit -1 HEAD || echo 'git') | tr -d '\n')

SEMVER=${VERSION}+${HEAD}-${BRANCH}
BUILD_TIME=$(shell date +%s)

CLIENT_ID="c19o8hor03fsa23cywutub8pu82ovo"

build:
	go build -tags netgo -ldflags "-X main.appVersion=${SEMVER} -X main.appBuildTime=${BUILD_TIME} -X main.appBuildUser=${USER} -X main.appClientId=${CLIENT_ID}"

fmt:
	go fmt `go list ./...`

vet:
	go vet `go list ./...`

check:
	go test `go list ./...`

deps:
	@echo "Depends on libvlc"
