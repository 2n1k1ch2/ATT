## Avito Test Task


собрать make docker-up

Есть доп.ручка по сборке статистики всех пользователей

Есть конфиг линтера, запускать make lint

Доп возможности описаны в make help

Есть сид с готовыми валидными данными чтобы потестить ручки.

Есть сваггер для удобства тестирования ручек (опционально pgadmin для удобного просмотра бд)

## Основные эндпоинты
Team 

POST /api/v1/team/add - Создание новой команды

GET /api/v1/team/get - Получение списка пользователей команды

Users

POST /api/v1/users/setIsActive - Изменение флага активности пользователя

GET /api/v1/users/getReview - Получение назначенных PR для пользователя

Pull Request 

POST /api/v1/pullRequest/create - Создание нового PR с автоматическим назначением ревьюверов

POST /api/v1/pullRequest/merge - Слияние PR (идемпотентная операция)

POST /api/v1/pullRequest/reassign - Переназначение ревьювера

Statistic

GET /api/v1/statistic/user - Статистика по пользователям

## Данные для pgadmin
Для удобного просмотра базы данных доступен PGAdmin:

URL: http://localhost:8081

Email: admin@avito.ru

Пароль: admin

Данные для подключения к БД в PGAdmin:
Host: db

Database: avito

Username: avito

Password: avito


### AfterDeadline
P.S Дополнительные задания которые я не успел сделать до 23.11.2025 Возможно будут лежать в бренче AfterDeadline (Если у вас будет желание посмотреть)
