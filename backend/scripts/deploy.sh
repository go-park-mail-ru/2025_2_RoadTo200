#!/usr/bin/bash

HOST=ubuntu.vk # Хост ОС. Пример: user@ххх.хх.хх.ххх
DOCS=false     # Флаг сборки документации
CONF=false     # Флаг отправки конфигурации

# Запускать из ./backend

if $DOCS; then
  make build-docs
fi

make build

if $CONF; then
  scp ./config/config.yaml $HOST:/home/ubuntu/app/back/config/
fi
scp ./../build/main $HOST:/home/ubuntu/app/back/build/
