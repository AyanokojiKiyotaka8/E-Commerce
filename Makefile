GOBIN := $(shell go env GOBIN)
ifeq ($(GOBIN),)
GOBIN := $(shell go env GOPATH)/bin
endif

PROTOC_GEN_GO := $(GOBIN)/protoc-gen-go
PROTOC_GEN_GRPC := $(GOBIN)/protoc-gen-go-grpc
PROTOC_GEN_GATEWAY := $(GOBIN)/protoc-gen-grpc-gateway

check-plugins:
	@command -v $(PROTOC_GEN_GO) >/dev/null || (echo "Missing protoc-gen-go" && exit 1)
	@command -v $(PROTOC_GEN_GRPC) >/dev/null || (echo "Missing protoc-gen-go-grpc" && exit 1)
	@command -v $(PROTOC_GEN_GATEWAY) >/dev/null || (echo "Missing protoc-gen-grpc-gateway" && exit 1)

pproto: check-plugins
	@protoc \
		-I . \
		-I ./googleapis \
		--plugin=protoc-gen-go=$(PROTOC_GEN_GO) \
		--plugin=protoc-gen-go-grpc=$(PROTOC_GEN_GRPC) \
		--plugin=protoc-gen-grpc-gateway=$(PROTOC_GEN_GATEWAY) \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
		product_service/proto/product.proto

gateway: 
	@go build -o bin/gateway ./gateway
	@./bin/gateway

product:
	@go build -o bin/product_service ./product_service/cmd
	@./bin/product_service

.PHONY: check-plugins pproto gateway product
