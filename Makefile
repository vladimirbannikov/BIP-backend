# TODO: Add tests

LOCAL_BIN:=$(CURDIR)/bin
DEPLOY_DIR:=$(CURDIR)/deployments

# run in docker
.PHONY: run-in-docker
run-in-docker: build
	cd $(DEPLOY_DIR) && docker-compose up --build

# build app
.PHONY: build
build:
	go mod download && CGO_ENABLED=0  go build \
		-o ./bin/main$(shell go env GOEXE) ./cmd/main.go
