#!/bin/bash
jq -r '.code' /playground/app/suite.json | base64 -d > /playground/app/main.go
cd /playground/app || exit
go mod tidy
timeout 60 go build -o main .
jq -r '.input' /playground/app/suite.json | xargs timeout -v "$TIMEOUT" ./main
