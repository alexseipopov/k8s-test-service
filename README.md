# T-Ticker - Программа для периодической записи в PostgreSQL

Программа на Go, которая каждые 5 секунд добавляет новую запись в базу данных PostgreSQL с использованием sqlx.

## Возможности

- Автоматическое создание таблицы в PostgreSQL при первом запуске
- Добавление записей каждые 5 секунд
- Отображение прогресса в консоли
- Подсчет общего количества записей
- Показ последних записей каждые 5 итераций
- Обработка ошибок
- Поддержка переменных окружения для подключения к БД

## Технологии

- **Go** - основной язык программирования
- **PostgreSQL** - база данных
- **sqlx** - расширенная библиотека для работы с SQL
- **Docker Compose** - для локального запуска PostgreSQL

## Установка и запуск

### 1. Запуск PostgreSQL через Docker Compose

```bash
docker-compose up -d
```

Это запустит PostgreSQL на порту 5432 с следующими параметрами:
- База данных: `ticker_db`
- Пользователь: `ticker_user`
- Пароль: `ticker_password`

### 2. Установка зависимостей Go

```bash
go mod tidy
```

### 3. Запуск программы

```bash
go run main.go
```

### 4. Остановка

Для остановки программы нажмите `Ctrl+C`

Для остановки PostgreSQL:
```bash
docker-compose down
```

## Переменные окружения

Программа поддерживает переменную окружения `DATABASE_URL` для настройки подключения к БД:

```bash
export DATABASE_URL="host=localhost port=5432 user=ticker_user password=ticker_password dbname=ticker_db sslmode=disable"
go run main.go
```

Если переменная не задана, используются значения по умолчанию.

## Структура БД

Программа создает таблицу `records` со следующими полями:
- `id` - SERIAL PRIMARY KEY (автоинкрементный)
- `message` - TEXT NOT NULL (текстовое сообщение)
- `timestamp` - TIMESTAMP WITH TIME ZONE (дата и время создания записи)

## Файлы проекта

- `main.go` - основной код программы
- `go.mod` - файл зависимостей Go
- `docker-compose.yml` - конфигурация PostgreSQL
- `.gitignore` - файл исключений для Git
- `README.md` - документация

## Пример вывода

```
Программа запущена. Записи будут добавляться каждые 5 секунд...
Подключение к БД: host=localhost port=5432 user=ticker_user password=ticker_password dbname=ticker_db sslmode=disable
Нажмите Ctrl+C для остановки
✓ Добавлена запись: Запись #1 - 2024-01-15 14:30:00 (всего записей: 1)
✓ Добавлена запись: Запись #2 - 2024-01-15 14:30:05 (всего записей: 2)
✓ Добавлена запись: Запись #3 - 2024-01-15 14:30:10 (всего записей: 3)
✓ Добавлена запись: Запись #4 - 2024-01-15 14:30:15 (всего записей: 4)
✓ Добавлена запись: Запись #5 - 2024-01-15 14:30:20 (всего записей: 5)
Последние записи:
  - ID: 5, Сообщение: Запись #5 - 2024-01-15 14:30:20, Время: 2024-01-15 14:30:20
  - ID: 4, Сообщение: Запись #4 - 2024-01-15 14:30:15, Время: 2024-01-15 14:30:15
  - ID: 3, Сообщение: Запись #3 - 2024-01-15 14:30:10, Время: 2024-01-15 14:30:10
```

## Проверка данных в БД

Для проверки данных в PostgreSQL можно использовать:

```bash
# Подключение к БД через Docker
docker exec -it t-ticker-postgres psql -U ticker_user -d ticker_db

# Просмотр всех записей
SELECT * FROM records ORDER BY timestamp DESC;

# Подсчет записей
SELECT COUNT(*) FROM records;
```

## Деплой в Kubernetes

### Предварительные требования

- Docker
- kubectl (настроенный для подключения к кластеру)
- Helm 3.x
- Доступ к Docker Hub (для публикации образов)

### Быстрый старт

1. **Сборка и публикация Docker образа**:
   ```bash
   # Сборка и отправка в Docker Hub
   ./scripts/build-and-push.sh v1.0.0
   
   # Или с тегом latest
   ./scripts/build-and-push.sh
   ```

2. **Деплой в Kubernetes**:
   ```bash
   # Деплой в namespace по умолчанию
   ./scripts/deploy.sh
   
   # Деплой в конкретный namespace
   ./scripts/deploy.sh my-namespace v1.0.0
   ```

3. **Проверка деплоя**:
   ```bash
   # Просмотр подов
   kubectl get pods -l app.kubernetes.io/name=t-ticker
   
   # Просмотр логов
   kubectl logs -l app.kubernetes.io/name=t-ticker -f
   
   # Просмотр всех ресурсов
   kubectl get all -l app.kubernetes.io/name=t-ticker
   ```

### Ручной деплой

Если вы предпочитаете ручной контроль:

1. **Сборка образа**:
   ```bash
   docker build -t alexseipopov/t-ticker:latest .
   docker push alexseipopov/t-ticker:latest
   ```

2. **Обновление зависимостей Helm**:
   ```bash
   helm dependency update ./helm/t-ticker
   ```

3. **Деплой с Helm**:
   ```bash
   helm upgrade --install t-ticker ./helm/t-ticker \
     --set image.tag=latest \
     --set postgresql.enabled=true \
     --wait
   ```

### Конфигурация

Основные параметры можно настроить в `helm/t-ticker/values.yaml`:

- **Ресурсы**: CPU и память для приложения и PostgreSQL
- **Реплики**: Количество экземпляров приложения
- **PostgreSQL**: Настройки базы данных
- **Ingress**: Настройки входящего трафика (если нужен)

### Мониторинг и логи

```bash
# Просмотр логов приложения
kubectl logs -l app.kubernetes.io/name=t-ticker -f

# Просмотр логов PostgreSQL
kubectl logs -l app.kubernetes.io/name=postgresql -f

# Просмотр событий
kubectl get events --sort-by=.metadata.creationTimestamp

# Описание пода для диагностики
kubectl describe pod -l app.kubernetes.io/name=t-ticker
```

### Подключение к PostgreSQL в кластере

```bash
# Порт-форвардинг для локального подключения
kubectl port-forward svc/t-ticker-postgresql 5432:5432

# Подключение через psql
psql -h localhost -p 5432 -U ticker_user -d ticker_db

# Или через kubectl exec
kubectl exec -it deployment/t-ticker-postgresql -- psql -U ticker_user -d ticker_db
```

### Очистка

```bash
# Удаление деплоя
./scripts/cleanup.sh

# Или ручное удаление
helm uninstall t-ticker
kubectl delete pvc --all  # Если нужно удалить данные БД
```

### Структура Helm чарта

```
helm/t-ticker/
├── Chart.yaml              # Метаданные чарта
├── values.yaml             # Значения по умолчанию
├── requirements.yaml       # Зависимости (PostgreSQL)
└── templates/
    ├── deployment.yaml     # Deployment приложения
    ├── service.yaml        # Service
    ├── serviceaccount.yaml # Service Account
    ├── ingress.yaml        # Ingress (опционально)
    ├── hpa.yaml           # Horizontal Pod Autoscaler
    └── _helpers.tpl       # Вспомогательные шаблоны
```

### Переменные окружения

Приложение использует следующие переменные окружения:

- `DATABASE_URL` - строка подключения к PostgreSQL (автоматически настраивается через Helm)

### Безопасность

- Приложение запускается под непривилегированным пользователем
- Используется read-only файловая система
- Настроены security contexts
- Service Account создается автоматически
