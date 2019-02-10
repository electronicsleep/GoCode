#!/bin/bash
set -e
go test -v
GOOS=linux go build gocode.go
set +e
docker rm gocode
set -e
docker build -t gocode .
docker run -p 8000:8000 --name gocode -i -t gocode
