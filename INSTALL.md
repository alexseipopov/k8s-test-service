# Установка и настройка для деплоя в Kubernetes

## Предварительные требования

### 1. Docker
```bash
# macOS (через Homebrew)
brew install docker

# Или скачать Docker Desktop
# https://www.docker.com/products/docker-desktop

# Проверка установки
docker --version
```

### 2. kubectl
```bash
# macOS (через Homebrew)
brew install kubectl

# Или через curl
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/darwin/amd64/kubectl"
chmod +x kubectl
sudo mv kubectl /usr/local/bin/

# Проверка установки
kubectl version --client
```

### 3. Helm
```bash
# macOS (через Homebrew)
brew install helm

# Или через curl
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

# Проверка установки
helm version
```

### 4. Настройка kubectl для подключения к кластеру

#### Для локального кластера (minikube):
```bash
# Установка minikube
brew install minikube

# Запуск кластера
minikube start

# Настройка kubectl
kubectl config use-context minikube
```

#### Для Docker Desktop:
```bash
# Включить Kubernetes в Docker Desktop
# Settings -> Kubernetes -> Enable Kubernetes

# Проверка контекста
kubectl config current-context
```

#### Для облачного кластера:
```bash
# AWS EKS
aws eks update-kubeconfig --region us-west-2 --name my-cluster

# Google GKE
gcloud container clusters get-credentials my-cluster --zone us-central1-a

# Azure AKS
az aks get-credentials --resource-group myResourceGroup --name myAKSCluster
```

## Проверка готовности

```bash
# Проверка подключения к кластеру
kubectl cluster-info

# Проверка узлов
kubectl get nodes

# Проверка Helm
helm version

# Проверка Docker
docker info
```

## Первый деплой

1. **Клонируйте репозиторий**:
   ```bash
   git clone <your-repo-url>
   cd t-ticker
   ```

2. **Соберите и отправьте образ**:
   ```bash
   ./scripts/build-and-push.sh v1.0.0
   ```

3. **Задеплойте в Kubernetes**:
   ```bash
   ./scripts/deploy.sh
   ```

4. **Проверьте статус**:
   ```bash
   kubectl get pods -l app.kubernetes.io/name=t-ticker
   kubectl logs -l app.kubernetes.io/name=t-ticker -f
   ```

## Возможные проблемы

### Docker не запущен
```bash
# Запуск Docker Desktop или Docker daemon
# На macOS: открыть Docker Desktop
# На Linux: sudo systemctl start docker
```

### kubectl не может подключиться к кластеру
```bash
# Проверка конфигурации
kubectl config view

# Проверка контекста
kubectl config current-context

# Список контекстов
kubectl config get-contexts
```

### Helm не может найти чарт
```bash
# Обновление зависимостей
helm dependency update ./helm/t-ticker

# Проверка чарта
helm lint ./helm/t-ticker
```

### Проблемы с правами доступа
```bash
# Сделать скрипты исполняемыми
chmod +x scripts/*.sh

# Проверить права
ls -la scripts/
```

## Полезные команды для отладки

```bash
# Просмотр событий в кластере
kubectl get events --sort-by=.metadata.creationTimestamp

# Описание пода
kubectl describe pod <pod-name>

# Просмотр логов
kubectl logs <pod-name> --previous

# Подключение к поду
kubectl exec -it <pod-name> -- /bin/sh

# Просмотр ресурсов
kubectl top pods
kubectl top nodes
```

## Следующие шаги

После успешного деплоя:

1. Настройте мониторинг (Prometheus, Grafana)
2. Настройте логирование (ELK Stack, Fluentd)
3. Настройте CI/CD pipeline
4. Добавьте health checks и readiness probes
5. Настройте autoscaling
6. Добавьте backup для PostgreSQL

