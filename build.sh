#!/bin/bash

export CGO_ENABLED=0
go build -o ./bin/loginservice ./cmd/authserver.go
