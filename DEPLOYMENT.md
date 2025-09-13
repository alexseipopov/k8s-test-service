# Команды для деплоя T-Ticker в Kubernetes

## 🚀 Быстрый старт

### 1. Сборка и публикация образа
```bash
# Сборка и отправка в Docker Hub с тегом latest
./scripts/build-and-push.sh

# Сборка с конкретным тегом
./scripts/build-and-push.sh v1.0.0
```

### 2. Деплой в Kubernetes
```bash
# Деплой в namespace по умолчанию
./scripts/deploy.sh

# Деплой в конкретный namespace
./scripts/deploy.sh production v1.0.0
```

### 3. Проверка статуса
```bash
# Просмотр подов
kubectl get pods -l app.kubernetes.io/name=t-ticker

# Просмотр логов
kubectl logs -l app.kubernetes.io/name=t-ticker -f

# Просмотр всех ресурсов
kubectl get all -l app.kubernetes.io/name=t-ticker
```

## 🔧 Ручные команды

### Docker
```bash
# Сборка образа
docker build -t alexseipopov/t-ticker:latest .

# Отправка в Docker Hub
docker push alexseipopov/t-ticker:latest

# Локальный тест
docker run --rm alexseipopov/t-ticker:latest
```

### Helm
```bash
# Обновление зависимостей
helm dependency update ./helm/t-ticker

# Деплой
helm upgrade --install t-ticker ./helm/t-ticker \
  --set image.tag=latest \
  --set postgresql.enabled=true \
  --wait

# Просмотр статуса
helm status t-ticker

# Удаление
helm uninstall t-ticker
```

### Kubernetes
```bash
# Просмотр подов
kubectl get pods

# Просмотр сервисов
kubectl get services

# Просмотр логов
kubectl logs deployment/t-ticker

# Порт-форвардинг для PostgreSQL
kubectl port-forward svc/t-ticker-postgresql 5432:5432

# Подключение к PostgreSQL
kubectl exec -it deployment/t-ticker-postgresql -- psql -U ticker_user -d ticker_db
```

## 🧹 Очистка

```bash
# Автоматическая очистка
./scripts/cleanup.sh

# Ручная очистка
helm uninstall t-ticker
kubectl delete pvc --all
```

## 📊 Мониторинг

```bash
# Просмотр событий
kubectl get events --sort-by=.metadata.creationTimestamp

# Описание пода
kubectl describe pod -l app.kubernetes.io/name=t-ticker

# Просмотр ресурсов
kubectl top pods -l app.kubernetes.io/name=t-ticker

# Просмотр логов PostgreSQL
kubectl logs -l app.kubernetes.io/name=postgresql -f
```

## 🔍 Отладка

```bash
# Проверка подключения к кластеру
kubectl cluster-info

# Проверка конфигурации Helm
helm template t-ticker ./helm/t-ticker

# Проверка манифестов
helm get manifest t-ticker

# Просмотр значений
helm get values t-ticker
```

## 📝 Полезные команды

```bash
# Просмотр всех ресурсов в namespace
kubectl get all

# Просмотр секретов
kubectl get secrets

# Просмотр конфигурационных карт
kubectl get configmaps

# Просмотр PVC
kubectl get pvc

# Просмотр ingress
kubectl get ingress
```

