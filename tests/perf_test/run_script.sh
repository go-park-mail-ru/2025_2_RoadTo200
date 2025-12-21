#!/bin/bash

set -e

# Конфигурация
API_URL=${2:-"https://terabithia.online"}
WRITE_URL="register"
READ_URL="login"
RESULTS_DIR="results"
mkdir -p "$RESULTS_DIR"

echo "Цель: $API_URL"
echo ""

test_create() {
    echo "📝 Тестирование создания 100k сущностей..."

    wrk -t10 -c1000 -d100s -s register.lua --latency "${API_URL}/api/${WRITE_URL}" > "${RESULTS_DIR}/create_report.txt"

    echo "✅ Создание завершено. RPS: 100, Длительность: 100 с"
    cat "${RESULTS_DIR}/create_report.txt"
}

test_read() {
    echo "👁️ Тестирование чтения сущностей..."

    wrk -t10 -c1000 -d100s -s login.lua --latency "${API_URL}/api/${READ_URL}" > "${RESULTS_DIR}/read_report.txt"

    echo "✅ Чтение завершено. RPS: 100, Длительность: 5 с"
    cat "${RESULTS_DIR}/read_report.txt"
}

# Главный цикл
main() {
    test_create
    test_read
}

main
