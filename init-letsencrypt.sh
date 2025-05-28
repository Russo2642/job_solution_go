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

# Создаем директории с правильными правами
mkdir -p "$data_path/conf/live/$domains"
mkdir -p "$data_path/www"
chmod -R 755 "$data_path"

echo "Настройка временного SSL сертификата..."
path="$data_path/conf/live/$domains"

# Создаем директорию, если её нет
mkdir -p $path

# Запускаем docker-compose с правильными переменными окружения
docker-compose run --rm --entrypoint "\
  openssl req -x509 -nodes -newkey rsa:$rsa_key_size -days 1\
    -keyout '$path/privkey.pem' \
    -out '$path/fullchain.pem' \
    -subj '/CN=localhost'" certbot

# Проверяем, созданы ли файлы сертификатов
if [ ! -f "$path/privkey.pem" ] || [ ! -f "$path/fullchain.pem" ]; then
  echo "Ошибка: не удалось создать временные сертификаты."
  exit 1
fi

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

echo "Настройка постоянного обновления сертификатов..."
# Создаем файл для cron, если он еще не существует
if [ ! -f "/etc/cron.d/certbot-renewal" ]; then
  echo "0 0,12 * * * root docker-compose -f $(pwd)/docker-compose.yml run --rm certbot renew" | sudo tee /etc/cron.d/certbot-renewal > /dev/null
fi

echo "Перезапуск nginx с новыми сертификатами..."
docker-compose restart nginx

echo "Готово! Теперь ваш сайт должен быть доступен по HTTPS."
