## Avito Test Task


собрать make docker-up

Есть доп.ручка по сборке статистики всех пользователей

Есть конфиг линтера, запускать linter make lint

Доп возможности описаны в make help

Есть сид с готовыми валидными данными чтобы потестить ручки.

Есть сваггер для удобства тестирования ручек (опционально pgadmin для удобного просмотра бд)

## Основные эндпоинты
POST /segments - создание сегмента

DELETE /segments - удаление сегмента

POST /user-segments - управление сегментами пользователя

GET /user-segments/{id} - получение сегментов пользователя

GET /stats - статистика по всем пользователям
## Данные для pgadmin
PGAdmin
Для удобного просмотра базы данных доступен PGAdmin:

URL: http://localhost:5050

Email: admin@avito.ru

Пароль: admin

Данные для подключения к БД в PGAdmin:
Host: db

Database: avito

Username: avito

Password: avito


### AfterDeadline
P.S Дополнительные задания которые я не успел сделать до 23.11.2025 Возможно будут лежать в бренче AfterDeadline (Если у вас будет желание посмотреть)
