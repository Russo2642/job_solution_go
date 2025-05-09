#!/bin/bash

if ! [ -x "$(command -v docker-compose)" ]; then
  echo 'Ошибка: docker-compose не установлен.' >&2
  exit 1
fi

domains=(77.240.38.137)
rsa_key_size=4096
data_path="./certbot"
email="your-email@example.com" # Укажите ваш email

if [ -d "$data_path" ]; then
  read -p "Существующие данные в $data_path. Продолжить и перезаписать существующие сертификаты? (y/N) " decision
  if [ "$decision" != "Y" ] && [ "$decision" != "y" ]; then
    exit
  fi
fi

if [ ! -d "$data_path/conf/live/$domains" ]; then
  echo "Создание директорий для хранения сертификатов..."
  mkdir -p "$data_path/conf/live/$domains"
  mkdir -p "$data_path/www"
fi

echo "Настройка временного SSL сертификата..."
path="/etc/letsencrypt/live/$domains"
mkdir -p "$data_path/conf/live/$domains"

# Создаем самоподписанный сертификат для начальной загрузки
docker-compose run --rm --entrypoint "\
  openssl req -x509 -nodes -newkey rsa:$rsa_key_size -days 1\
    -keyout '$path/privkey.pem' \
    -out '$path/fullchain.pem' \
    -subj '/CN=localhost'" certbot

echo "Получение сертификата Let's Encrypt..."
# Выключаем nginx
docker-compose down

# Запускаем nginx
docker-compose up --force-recreate -d nginx

# Запрашиваем сертификат
docker-compose run --rm --entrypoint "\
  certbot certonly --webroot -w /var/www/certbot \
    --email $email \
    --agree-tos \
    --no-eff-email \
    -d $domains \
    --force-renewal" certbot

echo "Перезапускаем nginx..."
docker-compose up -d

echo "Готово!" 