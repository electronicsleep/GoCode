#!/bin/bash
set -e
go test -v
GOOS=linux go build gocode.go
