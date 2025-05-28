#!/bin/bash

if ! [ -x "$(command -v docker-compose)" ]; then
  echo 'Ошибка: docker-compose не установлен.' >&2
  exit 1
fi

domains=(jobsolution.kz)
rsa_key_size=4096
data_path="./certbot"
email="gothiq2302@gmail.com"

# Проверяем, доступен ли домен
echo "Проверка DNS для $domains..."
if ! ping -c 1 $domains > /dev/null 2>&1; then
  echo "Ошибка: Домен $domains не доступен."
  echo "Убедитесь, что DNS настроен правильно и указывает на этот сервер."
  exit 1
fi

if [ -d "$data_path" ]; then
  read -p "Существующие данные в $data_path. Продолжить и перезаписать существующие сертификаты? (y/N) " decision
  if [ "$decision" != "Y" ] && [ "$decision" != "y" ]; then
    exit
  fi
fi

# Создаем все необходимые директории с правильными правами
echo "Создание директорий для сертификатов..."
rm -rf "$data_path/conf/live"
mkdir -p "$data_path/conf/live/$domains"
mkdir -p "$data_path/www"
chmod -R 755 "$data_path"

echo "Настройка временного SSL сертификата..."
path="$data_path/conf/live/$domains"

# Проверяем создание директории
if [ ! -d "$path" ]; then
  echo "Ошибка: Не удалось создать директорию $path"
  echo "Создаю директорию вручную..."
  sudo mkdir -p "$path"
  sudo chmod -R 755 "$path"
fi

# Создаем временный сертификат
echo "Создание временного сертификата..."
docker-compose run --rm --entrypoint "\
  mkdir -p /etc/letsencrypt/live/$domains && \
  openssl req -x509 -nodes -newkey rsa:$rsa_key_size -days 1\
    -keyout '/etc/letsencrypt/live/$domains/privkey.pem' \
    -out '/etc/letsencrypt/live/$domains/fullchain.pem' \
    -subj '/CN=localhost'" certbot

# Копируем сертификаты из контейнера в локальную директорию
echo "Копирование сертификатов..."
docker-compose run --rm --entrypoint "\
  cp -L /etc/letsencrypt/live/$domains/privkey.pem /etc/letsencrypt/live/$domains/fullchain.pem /var/www/certbot/" certbot

# Перезапускаем nginx для применения временных сертификатов
echo "Перезапуск nginx..."
docker-compose restart nginx

# Даем nginx время на перезапуск
sleep 5

echo "Получение сертификата Let's Encrypt..."
docker-compose run --rm --entrypoint "\
  certbot certonly --webroot -w /var/www/certbot \
    --email $email \
    --agree-tos \
    --no-eff-email \
    -d $domains \
    --force-renewal" certbot

# Проверяем успешность получения сертификата
if [ $? -ne 0 ]; then
  echo "Ошибка при получении сертификата Let's Encrypt."
  echo "Проверьте настройки DNS и доступность вашего сервера из интернета."
  echo "Временно настройте nginx на работу без SSL."
  exit 1
fi

echo "Перезапуск nginx с новыми сертификатами..."
docker-compose restart nginx

echo "Готово! Теперь ваш сайт должен быть доступен по HTTPS."
