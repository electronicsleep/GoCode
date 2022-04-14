#!/bin/bash
set -e
go test -v
go build gocode.go
./gocode
