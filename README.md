# Library Service API

Этот проект реализует сервис для работы с онлайн-библиотекой песен, предоставляющий REST API для управления песнями и обогащения данных из внешнего API.

## Описание задания

### Функционал:
1. REST API:
    - Получение списка песен с фильтрацией и пагинацией.
    - Получение текста песни с пагинацией по куплетам.
    - Добавление новой песни.
    - Удаление песни.
    - Изменение данных песни.
2. Данные хранятся в PostgreSQL (с использованием миграций).
3. Конфигурационные данные вынесены в `.env`.
4. API задокументировано с помощью Swagger.

---

## Используемые технологии

- **Язык:** Go (Golang)
- **Фреймворк для документации:** [Swaggo](https://github.com/swaggo/swag)
- **PostgreSQL** в качестве хранилища данных
- **Миграции:** [golang-migrate](https://github.com/golang-migrate/migrate)
- **Логирование:** [Logrus](https://github.com/sirupsen/logrus)
- **Конфигурация:** [Viper](https://github.com/spf13/viper)
- **Docker:** для контейнеризации приложения и базы данных.

---

## Установка и запуск

### Требования
- Docker и Docker Compose
- Go (для локального запуска)
- PostgreSQL (для локального тестирования)

### Локальный запуск через Docker
1. Склонируйте репозиторий:
   ```bash
   git clone https://github.com/username/library-service.git
   cd library-service
   
2. Запустите сервис с помощью Docker Compose:
    ```bash
   docker-compose up -d --build
   
3. Приложение будет доступно по адресу:
    ```
   http://localhost:8080
   
4. Swagger-документация доступна по адресу:
    ```
   http://localhost:8080/swagger/index.html

## Тестирование API

### Примеры запросов для тестирования API:

## Добавление песни:
    
    curl -X POST http://localhost:8080/songs \
    -H "Content-Type: application/json" \
    -d '{
      "group": "Muse",
      "song": "Supermassive Black Hole"
    }`

## Получение списка песен:

    curl -X GET http://localhost:8080/songs?group=Muse&page=1&per_page=10

## Обновление данных песни:

    curl -X PUT http://localhost:8080/songs/1 \
    -H "Content-Type: application/json" \
    -d '{
    "group": "Muse",
    "song": "Starlight",
    "release_date": "2006-07-18",
    "text": "Far away, the ship has taken me...",
    "link": "https://example.com/starlight"
    }'

## Удаление песни:

    curl -X DELETE http://localhost:8080/songs/1

    

