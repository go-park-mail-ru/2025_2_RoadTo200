#!/bin/bash

set -e

# Конфигурация
API_URL=${2:-"https://terabithia.online"}
WRITE_URL="register"
READ_URL="login"
RESULTS_DIR="results"
RESPS=30
mkdir -p "$RESULTS_DIR"

echo "Цель: $API_URL"
echo ""

test_create() {
    wrk -t10 -c${RESPS} -d60s -s register.lua --latency "${API_URL}/api/${WRITE_URL}" > "${RESULTS_DIR}/create_report.txt"

    cat "${RESULTS_DIR}/create_report.txt"
}

test_read() {
    wrk -t10 -c${RESPS} -d60s -s login.lua --latency "${API_URL}/api/${READ_URL}" > "${RESULTS_DIR}/read_report.txt"

    cat "${RESULTS_DIR}/read_report.txt"
}

# Главный цикл
main() {
    test_create
    test_read
}

main
