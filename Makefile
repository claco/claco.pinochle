.DEFAULT_GOAL := help
.PHONY: all build bump clean container generate help lint test version

ROOT       := ${PWD}
GOARCH     := ${shell go env GOARCH}
GOOS       := ${shell go env GOOS}
BINARY     := pinochle
PACKAGE    := github.com/claco/claco.pinochle
DISTDIR    := ${ROOT}/dist
BINDIR     := ${DISTDIR}/bin
VENV       := .venv
PROTOS     := protos
DOCKERFILE := ./Dockerfile
VCSTAG     := ${shell git describe --always}
VCSREV     := ${shell git rev-parse --verify HEAD}
VCSTIME    := ${shell git show -s --format=%cd --date=iso-strict HEAD}
VERSION    := latest
DEBUG      :=
GCFLAGS    := $(if ${DEBUG},all=-N -l,)
LDFLAGS    := $(if ${DEBUG},,-s -w) \
	-X '${PACKAGE}/build.VcsTag=${or ${VCSTAG},}' \
	-X '${PACKAGE}/build.VcsRevision=${or ${VCSREV},}' \
	-X '${PACKAGE}/build.VcsTime=${or ${VCSTIME},}' \
	-X '${PACKAGE}/build.Version=${VERSION}'
BUILDARGS  := --load --pull --file=${DOCKERFILE} \
	--build-arg=DEBUG=${DEBUG} \
	--build-arg=VCSTAG=${VCSTAG} \
	--build-arg=VCSTIME=${VCSTIME} \
	--build-arg=VCSREV=${VCSREV} \
	--build-arg=VERSION=${VERSION}
TARGET     := ./...

export CGO_ENABLED := 0
export GO_EXTLINK_ENABLED := 0

${BINDIR}:
	@mkdir -p $@

${BINDIR}/gocover-cobertura: ${BINDIR}
	@GOBIN=${BINDIR} go install github.com/richardlt/gocover-cobertura@latest

${BINDIR}/golangci-lint: ${BINDIR}
	@GOBIN=${BINDIR} go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

${BINDIR}/gocyclo: ${BINDIR}
	@GOBIN=${BINDIR} go install github.com/fzipp/gocyclo/cmd/gocyclo@latest

${BINDIR}/pre-commit: ${BINDIR} ${VENV}
	@${VENV}/bin/pip install pre-commit &>/dev/null && cp ${VENV}/bin/pre-commit ${BINDIR}/pre-commit

${BINDIR}/protoc-gen-go: ${BINDIR}
	@GOBIN=${BINDIR} go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	@GOBIN=${BINDIR} go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@GOBIN=${BINDIR} go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@GOBIN=${BINDIR} go install github.com/go-swagger/go-swagger/cmd/swagger@latest

${VENV}:
	@python3 -m venv ${VENV} && .venv/bin/pip install --quiet --upgrade pip pre-commit

all: ## clean, build, test, and dist
all: clean build lint test


build:  ## build project
build: ${BINDIR}
	@GOOS=${GOOS} GOARCH=${GOARCH} go build --gcflags="${GCFLAGS}" --ldflags="${LDFLAGS}" -o ${BINDIR}/${BINARY} .


bump: ## bump mod versions and tidy up
bump:
	@go get -u=patch
	@go mod tidy


clean: ## clean project
clean:
	@${RM} -rf ${BINDIR} ${DISTDIR} ${VENV}
	@docker rmi $$(docker images -q ${BINARY}) &>/dev/null || true


container: ## build container
container:
	@docker buildx build ${BUILDARGS} --target=pinochle --tag=${BINARY}:${VERSION} .
	@docker buildx build ${BUILDARGS} --target=pinochle --platform=linux/arm64 --tag=${BINARY}:${VERSION}-linux-arm64 .
	@docker buildx build ${BUILDARGS} --target=pinochle --platform=linux/amd64 --tag=${BINARY}:${VERSION}-linux-amd64 .
ifdef DEBUG
	@docker buildx build ${BUILDARGS} --target=debug --tag=${BINARY}:${VERSION}-debug .
endif


generate: ## generate grpc files
generate: ${BINDIR}/protoc-gen-go
	@PATH=${BINDIR}:${PATH} protoc \
		--proto_path==${PROTOS} \
		--go_out=./pb \
		--go_opt=paths=source_relative \
		--go-grpc_out=./pb \
		--go-grpc_opt=paths=source_relative \
		--openapiv2_out=./pb \
		--openapiv2_opt=logtostderr=true \
	${PROTOS}/*.proto


help: ## display this help
help:
	@echo "Usage: make [target] [argument=value] ..."
	@echo
	@egrep "^(.+)\:\s+##\ (.+)" ${MAKEFILE_LIST} | column -t -c 2 -s ":#"
	@echo


lint: ## lint project
lint: ${BINDIR}/pre-commit ${BINDIR}/golangci-lint ${BINDIR}/gocyclo
	@PATH=${BINDIR}:${PATH} pre-commit run --all-files


test: ## run project tests
test: build ${BINDIR}/gocover-cobertura
	@go test --race --covermode=atomic --coverprofile=coverage.out --coverpkg=${TARGET} ${TARGET} \
		&& ${BINDIR}/gocover-cobertura < coverage.out > coverage.xml


version: ## show build version
version: build
	@${BINDIR}/${BINARY} --version
