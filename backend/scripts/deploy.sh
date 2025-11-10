#!/usr/bin/bash

HOST=ubuntu.vk # Хост ОС. Пример: user@ххх.хх.хх.ххх
DOCS=0     # Флаг сборки документации
CONF=0     # Флаг отправки конфигурации

# Запускать из ./backend

if [[ $DOCS -eq 1 ]]; then
  make build-docs &&
  scp ./docs/swagger.json $HOST:/home/ubuntu/app/back/docs/
fi

make build || exit 1

if [[ $CONF -eq 1 ]]; then
  scp ./config/config.yaml $HOST:/home/ubuntu/app/back/config/
fi
scp ./../build/main $HOST:/home/ubuntu/app/back/build/
