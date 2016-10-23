#!/usr/bin/env bash

set -e
set -u

go test -coverprofile=coverage.out
go tool cover -html=coverage.out
