#!/bin/bash

# Проверяем, что контейнер запущен
if ! docker ps | grep -q postgres; then
  echo "Контейнер PostgreSQL не запущен. Запустите его командой 'docker-compose up -d'."
  exit 1
fi

# Функция для выполнения SQL-команд в контейнере
exec_psql() {
  docker exec -i postgres psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} <<EOF
$1
EOF
}

# Загружаем переменные из .env
set -a
source .env
set +a

# Основные настройки
echo "Настройка PostgreSQL в контейнере..."

# 1. Создаем пользователя приложения
echo "Создание пользователя app_user..."
exec_psql "
CREATE USER app_user WITH PASSWORD '${APP_USER_PASSWORD:-App_Pass_!456}';
GRANT CONNECT ON DATABASE $POSTGRES_DB TO app_user;
"

# 2. Настраиваем права
echo "Настройка прав для app_user..."
exec_psql "
GRANT CREATE ON SCHEMA public TO biba;
GRANT USAGE ON SCHEMA public TO app_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO biba;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO app_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO app_user;
"

# 3. Безопасность: отзываем публичные права
echo "Настройка безопасности..."
exec_psql "
REVOKE ALL ON DATABASE $POSTGRES_DB FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM PUBLIC;
"

# 4. Включаем логирование (требует перезагрузки контейнера)
echo "Включение логирования..."
exec_psql "
ALTER SYSTEM SET log_statement = 'all';
ALTER SYSTEM SET log_connections = 'on';
ALTER SYSTEM SET log_disconnections = 'on';
"

# 5. Создаем тестовую таблицу (пример)
echo "Создание тестовой таблицы..."
exec_psql "
CREATE TABLE IF NOT EXISTS test_table (
  id SERIAL PRIMARY KEY,
  data VARCHAR(100)
);
GRANT ALL PRIVILEGES ON TABLE test_table TO app_user;
"

# Перезагружаем контейнер для применения настроек
echo "Перезагрузка контейнера для применения настроек..."
docker restart postgres

# Ждем 5 секунд пока PostgreSQL перезапустится
sleep 5

# Проверяем подключение новым пользователем
echo "Проверка подключения app_user..."
docker exec -it postgres psql -U app_user -d "$POSTGRES_DB" -c "SELECT 'Успешное подключение!' AS message;"

echo "
Настройка завершена!
Данные для подключения:
- Хост: localhost
- Порт: 40502
- База: $POSTGRES_DB
- Пользователь приложения: app_user
"