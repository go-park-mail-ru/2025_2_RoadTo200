#!/bin/bash

# Быстрый тест основных endpoints с cookie
BASE_URL="http://217.16.17.116:8080"
COOKIE_FILE="test_cookies.txt"

# Очистка старого cookie файла
rm -f "$COOKIE_FILE"

echo "=== Быстрый тест API с Cookie ==="

# 1. Регистрация
#echo "1. Регистрация"
#curl -s -X POST "$BASE_URL/api/register" \
#  -H "Content-Type: application/json" \
#  -c "$COOKIE_FILE" \
#  -d '{
#    "email": "quick_test@example.com",
#    "password": "password123",
#    "passwordConfirm": "password123"
#  }' | jq

# 2. Вход
echo -e "\n2. Вход"
curl -s -X POST "$BASE_URL/api/login" \
  -H "Content-Type: application/json" \
  -c "$COOKIE_FILE" \
  -d '{
    "email": "quick_test@example.com",
    "password": "password123"
  }' | jq

# 3. Проверка сессии
echo -e "\n3. Проверка сессии"
curl -s -X GET "$BASE_URL/api/session" \
  -b "$COOKIE_FILE"


# 4. Получение профиля
echo -e "\n4. Получение профиля"
curl -s -X GET "$BASE_URL/api/profile" \
  -b "$COOKIE_FILE" | jq

# 5. Получение ленты
echo -e "\n5. Получение ленты"
curl -s -X GET "$BASE_URL/api/feed?limit=5" \
  -b "$COOKIE_FILE" | jq

# 6. Получение мэтчей
echo -e "\n6. Получение мэтчей"
curl -s -X GET "$BASE_URL/api/match?limit=5" \
  -b "$COOKIE_FILE" | jq

# 7. Выход
echo -e "\n7. Выход"
curl -s -X POST "$COOKIE_FILE" \
  -b "$COOKIE_FILE" | jq

# Показать cookie
echo -e "\nСодержимое cookie файла:"
cat "$COOKIE_FILE"

# Очистка
rm -f "$COOKIE_FILE"

echo -e "\n=== Тест завершен ==="
