# 🚀 Полное руководство по CI/CD для T-Ticker

## 📋 **Обзор**

Этот проект настроен для автоматического деплоя в Kubernetes через GitHub Actions. Поддерживаются два окружения:
- **Staging** - для тестирования (ветка `develop`)
- **Production** - для продакшна (ветка `main`)

## 🏗️ **Архитектура CI/CD**

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Feature       │    │   Develop       │    │   Main          │
│   Branch        │───▶│   Branch        │───▶│   Branch        │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Pull Request  │    │   Staging       │    │   Production    │
│   Tests         │    │   Deployment    │    │   Deployment    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🚀 **Быстрый старт**

### 1. Настройка сервера

```bash
# Клонируйте репозиторий
git clone <your-repo-url>
cd t-ticker

# Настройте сервер (выберите один вариант)
./scripts/setup-server.sh staging
./scripts/setup-server.sh production
```

### 2. Настройка GitHub

1. **Добавьте секреты в GitHub** (см. [GITHUB_SETUP.md](GITHUB_SETUP.md))
2. **Создайте окружения** `staging` и `production`
3. **Настройте защиту веток**

### 3. Первый деплой

```bash
# Создайте develop ветку
git checkout -b develop
git push -u origin develop

# Создайте PR в main
# После merge произойдет автоматический деплой
```

## 📁 **Структура проекта**

```
t-ticker/
├── .github/workflows/          # GitHub Actions
│   └── ci-cd.yml              # Основной workflow
├── helm/t-ticker/             # Helm чарты
│   ├── templates/             # Kubernetes манифесты
│   ├── values.yaml            # Базовые значения
│   ├── values-staging.yaml    # Staging конфигурация
│   └── values-production.yaml # Production конфигурация
├── scripts/                   # Скрипты автоматизации
│   ├── setup-server.sh        # Настройка сервера
│   ├── deploy-manual.sh       # Ручной деплой
│   └── health-check.sh        # Проверка здоровья
├── main.go                    # Исходный код приложения
├── Dockerfile                 # Docker образ
└── README.md                  # Документация
```

## 🔄 **Workflow процесс**

### 1. Feature Development

```bash
# Создайте feature ветку
git checkout -b feature/new-feature

# Внесите изменения
git add .
git commit -m "feat: add new feature"
git push origin feature/new-feature

# Создайте PR в develop
```

### 2. Staging Deployment

```bash
# После merge в develop
git checkout develop
git pull origin develop

# Автоматически запустится:
# 1. Тесты
# 2. Сборка Docker образа
# 3. Деплой в staging
```

### 3. Production Deployment

```bash
# Создайте PR из develop в main
# После merge произойдет:
# 1. Тесты
# 2. Сборка Docker образа
# 3. Деплой в production
```

## 🛠️ **Команды для разработки**

### Локальная разработка

```bash
# Запуск локально
go run main.go

# Тесты
go test ./...

# Линтер
golangci-lint run

# Сборка
go build -o main .
```

### Docker

```bash
# Сборка образа
docker build -t alexseipopov/t-ticker:latest .

# Запуск локально
docker run --rm alexseipopov/t-ticker:latest

# Отправка в registry
docker push alexseipopov/t-ticker:latest
```

### Kubernetes

```bash
# Ручной деплой
./scripts/deploy-manual.sh staging v1.0.0
./scripts/deploy-manual.sh production v1.0.0

# Проверка здоровья
./scripts/health-check.sh staging
./scripts/health-check.sh production

# Просмотр логов
kubectl logs -n staging -l app.kubernetes.io/name=t-ticker -f
kubectl logs -n production -l app.kubernetes.io/name=t-ticker -f
```

## 🔧 **Конфигурация**

### Environment Variables

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `DATABASE_URL` | Строка подключения к PostgreSQL | Автоматически настраивается |
| `LOG_LEVEL` | Уровень логирования | `info` |
| `TICK_INTERVAL` | Интервал между записями | `5s` |

### Helm Values

Основные параметры в `values.yaml`:

```yaml
replicaCount: 1                    # Количество реплик
image:
  repository: alexseipopov/t-ticker
  tag: latest
resources:
  limits:
    cpu: 100m
    memory: 128Mi
postgresql:
  enabled: true                    # Включить PostgreSQL
  auth:
    password: "ticker_password"    # Пароль БД
```

## 📊 **Мониторинг**

### Логи

```bash
# Логи приложения
kubectl logs -n production -l app.kubernetes.io/name=t-ticker -f

# Логи PostgreSQL
kubectl logs -n production -l app.kubernetes.io/name=postgresql -f

# Все логи в namespace
kubectl logs -n production --all-containers=true -f
```

### Метрики

```bash
# Использование ресурсов
kubectl top pods -n production

# Статус подов
kubectl get pods -n production

# События
kubectl get events -n production --sort-by=.metadata.creationTimestamp
```

### Health Checks

```bash
# Автоматическая проверка
./scripts/health-check.sh production

# Ручная проверка
kubectl get pods -n production
kubectl get svc -n production
kubectl get ingress -n production
```

## 🔐 **Безопасность**

### Secrets Management

```bash
# Создание секрета
kubectl create secret generic app-secrets \
  --from-literal=database-password=secret-password \
  -n production

# Использование в Helm
helm upgrade t-ticker ./helm/t-ticker \
  --set postgresql.auth.password="$(kubectl get secret app-secrets -o jsonpath='{.data.database-password}' | base64 -d)"
```

### RBAC

```bash
# Создание ServiceAccount
kubectl create serviceaccount t-ticker-deployer -n production

# Создание ClusterRole
kubectl create clusterrole t-ticker-deployer \
  --verb=get,list,watch,create,update,patch,delete \
  --resource=pods,services,deployments

# Создание ClusterRoleBinding
kubectl create clusterrolebinding t-ticker-deployer \
  --clusterrole=t-ticker-deployer \
  --serviceaccount=production:t-ticker-deployer
```

## 🚨 **Troubleshooting**

### Частые проблемы

#### 1. Поды не запускаются

```bash
# Проверка событий
kubectl describe pod <pod-name> -n production

# Проверка логов
kubectl logs <pod-name> -n production

# Проверка ресурсов
kubectl top nodes
kubectl top pods -n production
```

#### 2. Проблемы с сетью

```bash
# Проверка сервисов
kubectl get svc -n production

# Проверка endpoints
kubectl get endpoints -n production

# Тест DNS
kubectl run test-dns --image=busybox --rm -it --restart=Never -- nslookup t-ticker-postgresql
```

#### 3. Проблемы с базой данных

```bash
# Проверка PostgreSQL
kubectl exec -it <postgres-pod> -n production -- psql -U ticker_user -d ticker_db

# Проверка подключения
kubectl exec -it <app-pod> -n production -- env | grep DATABASE_URL
```

#### 4. Проблемы с Ingress

```bash
# Проверка Ingress
kubectl get ingress -n production
kubectl describe ingress -n production

# Проверка Ingress Controller
kubectl get pods -n ingress-nginx
kubectl logs -n ingress-nginx -l app.kubernetes.io/name=ingress-nginx
```

### Rollback

```bash
# Откат Helm релиза
helm rollback t-ticker 1 -n production

# Откат к предыдущей версии образа
helm upgrade t-ticker ./helm/t-ticker \
  --set image.tag=previous-version \
  -n production
```

## 📚 **Полезные ссылки**

- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Helm Documentation](https://helm.sh/docs/)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Docker Documentation](https://docs.docker.com/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)

## 🆘 **Поддержка**

При возникновении проблем:

1. Проверьте логи GitHub Actions
2. Выполните health check
3. Проверьте статус ресурсов в кластере
4. Обратитесь к документации выше
5. Создайте issue в репозитории

