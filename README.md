# 📖 Цитатник (Quote Book)

Мини-сервис для хранения и управления цитатами с REST API интерфейсом

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org/)
[![Docker](https://img.shields.io/badge/Docker-✓-2496ED?logo=docker)](https://docker.com)

## 🚀 Возможности

- Добавление новых цитат
- Получение всех цитат
- Получение случайной цитаты
- Фильтрация цитат по автору
- Удаление цитат по ID

## 🛠️ Технологии

- **Язык**: Go (чистый stdlib + gorilla/mux)
- **Хранилище**: In-memory база (concurrent-safe)
- **Упаковка**: Docker
- **Логирование**: slog (структурированные логи)

## 📡 API Endpoints

| Метод | Путь | Описание |
|-------|------|----------|
| `POST` | `/quotes` | Добавить новую цитату |
| `GET` | `/quotes` | Получить все цитаты |
| `GET` | `/quotes/random` | Получить случайную цитату |
| `GET` | `/quotes?author={name}` | Фильтр по автору |
| `DELETE` | `/quotes/{id}` | Удалить цитату |

## 🏃 Запуск

### Требования
- Docker ([установка](https://docs.docker.com/get-docker/))

### Сборка и запуск
```bash
# Собрать образ
docker build -t quote_book .

# Запустить контейнер
docker run -p 8080:8080 quote_book
```

Сервис будет доступен по адресу:
http://localhost:8080

## 📝 Примеры запросов

### Добавить цитату
```bash
curl -X POST http://localhost:8080/quotes \
  -H "Content-Type: application/json" \
  -d '{"author":"Confucius", "quote":"Life is simple..."}'
```

### Получить все цитаты
```bash
curl http://localhost:8080/quotes
```

### Получить случайную цитату

```bash
curl http://localhost:8080/quotes/random
```

### Получить цитаты автора

```bash
curl "http://localhost:8080/quotes?author=Confucius"
```

### Удалить цитату

```bash
curl -X DELETE http://localhost:8080/quotes/1
```

## Особенности реализации

- **In-Memory база данных:**
  - Оптимизированное хранение с индексами
  - Фоновая сборка мусора (GC)
  - Минимальные блокировки при операциях

- **Оптимизации:**
  - Быстрое получение случайной цитаты
  - Эффективное управление памятью
  - Поддержка высоких нагрузок


📊 Пример ответа

```json
[
  {
    "id": 1,
    "author": "Confucius",
    "quote": "Life is simple..."
  },
  {
    "id": 2,
    "author": "Einstein",
    "quote": "Imagination is more important than knowledge."
  }
]
```

