#!/usr/bin/bash

HOST=ubuntu.vk # Хост ОС. Пример: user@ххх.хх.хх.ххх
DOCS=0     # Флаг сборки документации
CONF=0     # Флаг отправки конфигурации
MIGR=1

# Запускать из ./backend

if [[ $MIGR -eq 1 ]]; then
  scp ./migrations/ddl.sql $HOST:/home/ubuntu/app/back/migrations/
  scp ./migrations/dml.sql $HOST:/home/ubuntu/app/back/migrations/
  scp ./migrations/drop.sql $HOST:/home/ubuntu/app/back/migrations/
fi

if [[ $DOCS -eq 1 ]]; then
  make build-docs &&
  scp ./docs/swagger.json $HOST:/home/ubuntu/app/back/docs/
fi

make build || exit 1

if [[ $CONF -eq 1 ]]; then
  scp ./config/config.yaml $HOST:/home/ubuntu/app/back/config/
fi
scp ./../build/main $HOST:/home/ubuntu/app/back/build/
