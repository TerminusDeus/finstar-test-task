#!/bin/bash

protoc \
 -I. \
 -I$GOPATH/pkg/mod/github.com/infobloxopen/protoc-gen-gorm@v0.20.1/ \
 --go_out=:. \
 --gorm_out==module=proto:. \
 orm.proto

protoc \
  -I. \
  -I$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.3.0/ \
  -I$GOPATH/pkg/mod/github.com/googleapis/googleapis@v0.0.0-20210401225656-56fc6d43fed7/ \
  -I$GOPATH/pkg/mod/github.com/envoyproxy/protoc-gen-validate@v0.5.1/ \
  --go_out=module=github.com/TerminusDeus/finstar-test-task:. \
  --go-grpc_out=module=github.com/TerminusDeus/finstar-test-task:. \
  --grpc-gateway_out=logtostderr=true,generate_unbound_methods=true,module=github.com/TerminusDeus/finstar-test-task:. \
  --swagger_out=logtostderr=true:../docs/swagger \
  --validate_out=lang=go,module=github.com/TerminusDeus/finstar-test-task:. \
  service.proto
