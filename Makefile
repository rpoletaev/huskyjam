VERSION := $(shell git describe --tags --abbrev=0)
GIT_SHA := $(shell git rev-parse --short HEAD)
DATE := $(shell date +%s)
MODULE := github.com/rpoletaev/huskyjam
GO = GO111MODULE=on CGO_ENABLED=0 go
GO_FLAGS := -mod=vendor -tags production -installsuffix cgo


LDFLAGS += -X $(MODULE)/pkg.Timestamp=$(DATE)
LDFLAGS += -X $(MODULE)/pkg.Version=$(VERSION)
LDFLAGS += -X $(MODULE)/pkg.GitSHA=$(GIT_SHA)

CMD_NAMES := $(foreach pb, $(wildcard ./cmd/*), $(pb)_cmd)

build: $(CMD_NAMES)

$(CMD_NAMES): ./cmd/%_cmd:
	@echo $*
	$(GO) build $(GO_FLAGS) -o bin/$* -ldflags "$(LDFLAGS)" ./cmd/$*

test::
	GO111MODULE=on go test -mod vendor -v ./cmd/clients

clean:
	@rm bin/*

.PHONY: gen tools
gen:
	wire ./cmd/service
	swag init -d ./http -g api.go -o ./cmd/service/docs
	mockgen -source=./internal/backend.go -destination=./mock/service_backend.go -package=mock
	mockgen -source=./http/accounts.go -destination=./mock/http_accounts.go -package=mock
	mockgen -source=./pkg/auth/tokens.go -destination=./mock/auth_tokens.go -package=mock

tools:
	GO111MODULE=off go get -u github.com/google/wire/...
	GO111MODULE=off go get -u github.com/golang/mock/mockgen