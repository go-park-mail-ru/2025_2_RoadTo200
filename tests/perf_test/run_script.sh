#!/bin/bash
# perf_test/run_all.sh - Полный цикл с Vegeta

set -e

# Конфигурация
ITERATION=${1:-1}
API_URL=${2:-"http://localhost:8080"}
RESULTS_DIR="results/iteration_$ITERATION"
mkdir -p "$RESULTS_DIR"

echo "🎯 Vegeta нагрузочное тестирование - Итерация $ITERATION"
echo "Цель: $API_URL"
echo ""

# 1. Сброс БД к исходной схеме
reset_db() {
    echo "🔄 Сброс БД..."
    psql -U admin -d Tinder -p 5431 -f init.sql
}

# 2. Генератор данных для создания сущностей
generate_post_data() {
    # Создаем 100k уникальных запросов для POST
    for i in $(seq 1 100000); do
        cat <<EOF
POST ${API_URL}/api/user
Content-Type: application/json

{
    "username": "user_${i}_$(date +%s)",
    "email": "user_${i}@test.com",
    "created_at": "$(date -Iseconds)"
}
EOF
    done
}

# 3. Тестирование создания (POST)
test_create() {
    echo "📝 Тестирование создания 100k сущностей..."

    # Вариант A: Постепенная нагрузка (рекомендуется)
    echo "POST ${API_URL}/api/users" > "${RESULTS_DIR}/targets.txt"
    echo "Content-Type: application/json" >> "${RESULTS_DIR}/targets.txt"

    # Создаем файл с 1000 записей для теста
    for i in $(seq 1 1000); do
        echo >> "${RESULTS_DIR}/targets.txt"
        echo "{\"username\":\"testuser_$i\",\"email\":\"user$i@test.com\"}" >> "${RESULTS_DIR}/targets.txt"
    done

    # Запускаем нагрузку: 50 RPS в течение 30 минут для создания 100k
    echo "POST ${API_URL}/api/users" | \
    vegeta attack \
        -rate=50 \                    # 50 запросов в секунду
        -duration=30m \               # 30 минут = 90k запросов
        -body="${RESULTS_DIR}/targets.txt" \
        -header="Content-Type: application/json" \
        > "${RESULTS_DIR}/create.bin"

    # Анализируем результаты
    vegeta report "${RESULTS_DIR}/create.bin" > "${RESULTS_DIR}/create_report.txt"
    vegeta report -type=json "${RESULTS_DIR}/create.bin" > "${RESULTS_DIR}/create_report.json"

    echo "✅ Создание завершено. RPS: 50, Длительность: 30 мин"
    tail -10 "${RESULTS_DIR}/create_report.txt"
}

# 4. Тестирование чтения (GET)
test_read() {
    echo "👁️ Тестирование чтения сущностей..."

    # Генерируем URL для чтения случайных пользователей (1-100000)
    for i in $(seq 1 1000); do
        echo "GET ${API_URL}/api/users/$((RANDOM % 100000 + 1))"
    done > "${RESULTS_DIR}/read_targets.txt"

    # Запускаем нагрузку: 100 RPS, 5 минут
    vegeta attack \
        -rate=100 \
        -duration=5m \
        -targets="${RESULTS_DIR}/read_targets.txt" \
        > "${RESULTS_DIR}/read.bin"

    vegeta report "${RESULTS_DIR}/read.bin" > "${RESULTS_DIR}/read_report.txt"
    vegeta report -type=json "${RESULTS_DIR}/read.bin" > "${RESULTS_DIR}/read_report.json"

    echo "✅ Чтение завершено. RPS: 100, Длительность: 5 мин"
    tail -10 "${RESULTS_DIR}/read_report.txt"
}

# 5. Анализ и сравнение
analyze() {
    echo "📊 Сводка итерации $ITERATION:"
    echo "=============================="

    echo "СОЗДАНИЕ (POST):"
    grep -E "Requests|Success|Latencies" "${RESULTS_DIR}/create_report.txt"

    echo ""
    echo "ЧТЕНИЕ (GET):"
    grep -E "Requests|Success|Latencies" "${RESULTS_DIR}/read_report.txt"

    # Сохраняем ключевые метрики для сравнения
    {
        echo "Итерация $ITERATION: $(date)"
        echo "POST - p95: $(grep "95%" "${RESULTS_DIR}/create_report.txt" | awk '{print $2}')"
        echo "POST - RPS: $(grep "Requests" "${RESULTS_DIR}/create_report.txt" | awk '{print $2}')"
        echo "GET - p95: $(grep "95%" "${RESULTS_DIR}/read_report.txt" | awk '{print $2}')"
        echo "GET - RPS: $(grep "Requests" "${RESULTS_DIR}/read_report.txt" | awk '{print $2}')"
    } > "${RESULTS_DIR}/summary.txt"

    cat "${RESULTS_DIR}/summary.txt"
}

# 6. Визуализация (опционально)
visualize() {
    if command -v vegeta &> /dev/null; then
        echo "📈 Генерация графиков..."
        vegeta plot "${RESULTS_DIR}/create.bin" > "${RESULTS_DIR}/create_plot.html"
        vegeta plot "${RESULTS_DIR}/read.bin" > "${RESULTS_DIR}/read_plot.html"
        echo "Откройте ${RESULTS_DIR}/*_plot.html в браузере"
    fi
}

# Главный цикл
main() {
    reset_db
    test_create
    test_read
    analyze
    visualize

    echo ""
    echo "✅ Итерация $ITERATION завершена!"
    echo "📁 Результаты в: ${RESULTS_DIR}/"
    echo ""
    echo "🔄 Для следующей итерации после оптимизаций:"
    echo "   ./run_all.sh $((ITERATION + 1))"
}

main
