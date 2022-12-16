.DEFAULT_GOAL := all
.PHONY: all build clean check configure dist distclean help install installdirs uinstall

DESTDIR := src/python/pinochle/grpc
DISTDIR := dist
LOGDIR	:= logs
VENV	:= .venv

$(shell mkdir -p ${LOGDIR})

help:
	@echo "Usage: make [target] [argument=value] ..."
	@echo
	@egrep "^(.+)\:\s+##\ (.+)" ${MAKEFILE_LIST} | column -t -c 2 -s ":#"
	@echo

all: ## clean, build, check, and dist
all: clean build check dist

build: ## build this project
build: configure ${VENV}/bin/coverage generate

check: ## run project checks/tests
check: build
	@${VENV}/bin/coverage run
	@${VENV}/bin/coverage xml
	@${VENV}/bin/coverage report

configure: ## configure this project
configure:
	@LOGDIR=${LOGDIR} ./configure 2>&1> ${LOGDIR}/configure.log

dist: ## make distributable
dist: build
	@${VENV}/bin/poetry build 2>&1> ${LOGDIR}/dist.log

distclean: ## reset project status
	@rm -rf ${DISTDIR} ${LOGDIR} ${VENV} .coverage* .pytest_cache

generate: ## regenerate the protobufs
generate:
	@${VENV}/bin/python -m grpc_tools.protoc \
			--proto_path=protos \
			--fatal_warnings \
			--include_imports \
			--descriptor_set_out="${DESTDIR}/pinochle.pb" \
			--grpc_python_out=${DESTDIR} \
			--openapiv2_out=${DESTDIR} \
			--python_out="${DESTDIR}" \
			--pyi_out="${DESTDIR}" \
		protos/*.proto

install: ## install pinochle
install: ${VENV}/bin/coverage

${VENV}/bin/activate: pyproject.toml
	@python3 -m venv "${VENV}" 2>&1> ${LOGDIR}/python.log
	@${VENV}/bin/pip install --upgrade pip poetry setuptools wheel 2>&1> ${LOGDIR}/pip.log

${VENV}/bin/coverage: ${VENV}/bin/activate pyproject.toml
	@${VENV}/bin/poetry install 2>&1> ${LOGDIR}/poetry.log
