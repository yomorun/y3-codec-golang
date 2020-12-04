#!/usr/bin/env bash

set -e

go tool pprof -http=":8081" http://localhost:6060/debug/cpu/profile
