#!/usr/bin/env bash

GO111MODULE=on go get github.com/golang/protobuf/protoc-gen-go google.golang.org/grpc/cmd/protoc-gen-go-grpc

go generate ./...
go build cmd/ag-diag.go
