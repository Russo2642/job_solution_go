#!/bin/bash

if ! [ -x "$(command -v docker-compose)" ]; then
  echo 'Ошибка: docker-compose не установлен.' >&2
  exit 1
fi

domains=(jobsolution.kz)
rsa_key_size=4096
data_path="./certbot"
email="gothiq2302@gmail.com"

if [ -d "$data_path" ]; then
  read -p "Существующие данные в $data_path. Продолжить и перезаписать существующие сертификаты? (y/N) " decision
  if [ "$decision" != "Y" ] && [ "$decision" != "y" ]; then
    exit
  fi
fi

mkdir -p "$data_path/conf/live/$domains"
mkdir -p "$data_path/www"

echo "Настройка временного SSL сертификата..."
path="$data_path/conf/live/$domains"

docker-compose run --rm --entrypoint "\
  openssl req -x509 -nodes -newkey rsa:$rsa_key_size -days 1\
    -keyout '$path/privkey.pem' \
    -out '$path/fullchain.pem' \
    -subj '/CN=localhost'" certbot

echo "Перезапускаем только nginx..."
docker-compose stop nginx
docker-compose up -d nginx

echo "Получение сертификата Let's Encrypt..."
docker-compose run --rm --entrypoint "\
  certbot certonly --webroot -w /var/www/certbot \
    --email $email \
    --agree-tos \
    --no-eff-email \
    -d $domains \
    --force-renewal" certbot

echo "Перезапускаем nginx после получения сертификата..."
docker-compose restart nginx

echo "Готово!"