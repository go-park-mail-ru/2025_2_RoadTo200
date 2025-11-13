#!/bin/bash

HOST=ubuntu.vk # Хост ОС. Пример: user@ххх.хх.хх.ххх
DOCS=0     # Флаг сборки документации
CONF=0     # Флаг отправки конфигурации
MIGR=0
BUILD=1

if [[ $BUILD -eq 1 ]]; then
  echo  "Building..."
  make build || exit 1

  echo "Deploy binary file"
  scp ./.build/main $HOST:/home/ubuntu/app/back/bin/ || echo "Error deploy binary"
fi

if [[ $MIGR -eq 1 ]]; then
  echo "Deploy migration files"
  scp -r ./migrations/* $HOST:/home/ubuntu/app/back/data/migrations/ || echo "Error deploy migrations"
fi

if [[ $DOCS -eq 1 ]]; then
  echo "Building swagger docs"
  make build-docs &&
  echo "Deploy swagger file"
  scp ./docs/swagger.json $HOST:/home/ubuntu/app/back/data/docs/ || echo "Error deploy swagger"
fi


if [[ $CONF -eq 1 ]]; then
  echo "Deploy config file"
  scp ./config/config.yaml $HOST:/home/ubuntu/app/back/config/ || echo "Error deploy config"
fi
