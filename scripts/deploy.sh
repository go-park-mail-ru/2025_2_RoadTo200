#!/bin/bash

HOST=$1
if [[ -z $HOST ]]; then
  HOST=ubuntu.vk # Хост ОС. Пример: user@ххх.хх.хх.ххх
fi
DOCS=1     # Флаг сборки документации
CONF=1     # Флаг отправки конфигурации
MIGR=1
INFR=1
BILD=1

if [[ $MIGR -eq 1 ]]; then
  echo "Deploy migration files"
  scp -o StrictHostKeyChecking=no -i ~/.ssh/deploy_key -r ./migrations/* $HOST:/home/ubuntu/app/back/data/migrations/ || echo "Error deploy migrations"
fi

if [[ $DOCS -eq 1 ]]; then
  echo "Building swagger docs"
  make build-docs &&
  echo "Deploy swagger file"
  scp -o StrictHostKeyChecking=no -i ~/.ssh/deploy_key ./api/auth/swagger.json $HOST:/home/ubuntu/app/back/data/docs/auth.json || echo "Error deploy auth swagger"
  scp -o StrictHostKeyChecking=no -i ~/.ssh/deploy_key ./api/core/swagger.json $HOST:/home/ubuntu/app/back/data/docs/core.json || echo "Error deploy core swagger"
  scp -o StrictHostKeyChecking=no -i ~/.ssh/deploy_key ./api/chat/swagger.json $HOST:/home/ubuntu/app/back/data/docs/chat.json || echo "Error deploy server swagger"
  scp -o StrictHostKeyChecking=no -i ~/.ssh/deploy_key ./api/server/swagger.json $HOST:/home/ubuntu/app/back/data/docs/server.json || echo "Error deploy server swagger"
fi


if [[ $CONF -eq 1 ]]; then
  echo "Deploy config file"
  scp -o StrictHostKeyChecking=no -i ~/.ssh/deploy_key ./config/* $HOST:/home/ubuntu/app/back/config/ || echo "Error deploy config"
fi

if [[ $INFR -eq 1 ]]; then
  echo "Deploy infra file"
  scp -o StrictHostKeyChecking=no -i ~/.ssh/deploy_key ./docker-compose.yml $HOST:/home/ubuntu/app/configs/docker-compose.yaml || echo "Error deploy docker config"
  scp -o StrictHostKeyChecking=no -i ~/.ssh/deploy_key -r ./infra/* $HOST:/home/ubuntu/app/configs/infra/ || echo "Error deploy metrics config"
fi

if [[ $BILD -eq 1 ]]; then
  echo  "Building..."
  make build-bin || exit 1

  echo "Deploy binary file"
  ssh -o StrictHostKeyChecking=no -i ~/.ssh/deploy_key $HOST sudo systemctl stop auth-trb.service chat-trb.service core-trb.service app-back.service
  scp -o StrictHostKeyChecking=no -i ~/.ssh/deploy_key ./.build/* $HOST:/home/ubuntu/app/back/bin/ || echo "Error deploy binary"
  ssh -o StrictHostKeyChecking=no -i ~/.ssh/deploy_key $HOST sudo systemctl start auth-trb.service chat-trb.service core-trb.service app-back.service
  ssh -o StrictHostKeyChecking=no -i ~/.ssh/deploy_key $HOST /home/ubuntu/app/back/app.sh status
fi
