#!/bin/bash
go test -v
GOOS=linux go build gocode.go
docker rm gocode
docker build -t gocode .
docker run -p 8000:8000 --name gocode -i -t gocode
