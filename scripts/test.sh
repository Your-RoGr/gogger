#!/bin/bash

cd "$(dirname $0)"

mkdir -p ./logs

go vet ../gogger/...
go test -race -cover -coverprofile=./logs/coverage.txt ../gogger/...
go tool cover -html ./logs/coverage.txt -o ./logs/coverage.html
go tool cover -func=./logs/coverage.txt
