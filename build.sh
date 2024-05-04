#!/bin/bash

# Build untuk arsitektur x86
env GOOS=linux GOARCH=386 go build -o andromodem_x86 ./cmd/andromodem

# Build untuk arsitektur amd64
env GOOS=linux GOARCH=amd64 go build -o andromodem_amd64 ./cmd/andromodem

# Build untuk arsitektur arm
env GOOS=linux GOARCH=arm go build -o andromodem_arm ./cmd/andromodem

# Build untuk arsitektur arm64
env GOOS=linux GOARCH=arm64 go build -o andromodem_arm64 ./cmd/andromodem
