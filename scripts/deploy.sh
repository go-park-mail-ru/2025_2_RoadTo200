#!/bin/bash

HOST=ubuntu.vk # Хост ОС. Пример: user@ххх.хх.хх.ххх
DOCS=0     # Флаг сборки документации
CONF=1     # Флаг отправки конфигурации
MIGR=0
BUILD=1

if [[ $MIGR -eq 1 ]]; then
  echo "Deploy migration files"
  scp -r ./migrations/* $HOST:/home/ubuntu/app/back/data/migrations/ || echo "Error deploy migrations"
fi

if [[ $DOCS -eq 1 ]]; then
  echo "Building swagger docs"
  make build-docs &&
  echo "Deploy swagger file"
  scp ./api/auth/swagger.json $HOST:/home/ubuntu/app/back/data/docs/auth.json || echo "Error deploy auth swagger"
  scp ./api/core/swagger.json $HOST:/home/ubuntu/app/back/data/docs/core.json || echo "Error deploy core swagger"
  scp ./api/chat/swagger.json $HOST:/home/ubuntu/app/back/data/docs/chat.json || echo "Error deploy server swagger"
  scp ./api/server/swagger.json $HOST:/home/ubuntu/app/back/data/docs/server.json || echo "Error deploy server swagger"
fi


if [[ $CONF -eq 1 ]]; then
  echo "Deploy config file"
  scp ./config/* $HOST:/home/ubuntu/app/back/config/ || echo "Error deploy config"
fi

if [[ $BUILD -eq 1 ]]; then
  echo  "Building..."
  make build-bin || exit 1

  echo "Deploy binary file"
  ssh ubuntu.vk sudo systemctl stop auth-trb.service chat-trb.service core-trb.service app-back.service
  scp ./.build/* $HOST:/home/ubuntu/app/back/bin/ || echo "Error deploy binary"
  ssh ubuntu.vk sudo systemctl start auth-trb.service chat-trb.service core-trb.service app-back.service
fi
