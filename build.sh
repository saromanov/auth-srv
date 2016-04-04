#!/usr/bin/env bash

echo "Get dependencies"
go get ./...
echo "Build service"
go build -o ./bin/service main.go
echo "Run service"
./bin/service