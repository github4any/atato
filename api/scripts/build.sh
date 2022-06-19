#!/bin/bash

go build -o mockredis-server cmd/server/main.go && go build -o mockredis-cli cmd/client/main.go