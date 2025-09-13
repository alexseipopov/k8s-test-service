# Настройка GitHub для CI/CD

## 🔐 **Настройка GitHub Secrets**

### 1. Перейдите в настройки репозитория

1. Откройте ваш репозиторий на GitHub
2. Перейдите в **Settings** → **Secrets and variables** → **Actions**
3. Нажмите **New repository secret**

### 2. Добавьте следующие секреты:

#### **Docker Hub секреты:**
```
DOCKER_USERNAME = alexseipopov
DOCKER_PASSWORD = <ваш-docker-hub-токен>
```

#### **Kubernetes конфигурации:**
```
KUBE_CONFIG_STAGING = <base64-encoded-kubeconfig-for-staging>
KUBE_CONFIG_PRODUCTION = <base64-encoded-kubeconfig-for-production>
```

### 3. Создание Docker Hub токена

1. Войдите в [Docker Hub](https://hub.docker.com)
2. Перейдите в **Account Settings** → **Security**
3. Нажмите **New Access Token**
4. Выберите **Read, Write, Delete** права
5. Скопируйте токен и добавьте в `DOCKER_PASSWORD`

### 4. Получение kubeconfig для CI/CD

#### Для staging окружения:
```bash
# На staging сервере
kubectl config view --raw --minify > kubeconfig-staging.yaml
base64 -w 0 kubeconfig-staging.yaml
```

#### Для production окружения:
```bash
# На production сервере
kubectl config view --raw --minify > kubeconfig-production.yaml
base64 -w 0 kubeconfig-production.yaml
```

## 🌍 **Настройка GitHub Environments**

### 1. Создание окружений

1. Перейдите в **Settings** → **Environments**
2. Создайте окружения:
   - `staging`
   - `production`

### 2. Настройка защиты окружений

#### Для staging:
- **Required reviewers**: 0 (автоматический деплой)
- **Wait timer**: 0 минут
- **Deployment branches**: `develop` branch only

#### Для production:
- **Required reviewers**: 1-2 человека
- **Wait timer**: 5 минут
- **Deployment branches**: `main` branch only

### 3. Добавление environment-specific секретов

#### Staging environment:
```
KUBE_CONFIG_STAGING = <staging-kubeconfig>
DOMAIN_STAGING = t-ticker-staging.yourdomain.com
```

#### Production environment:
```
KUBE_CONFIG_PRODUCTION = <production-kubeconfig>
DOMAIN_PRODUCTION = t-ticker.yourdomain.com
```

## 🔄 **Настройка веток и workflow**

### 1. Создание веток

```bash
# Создание develop ветки
git checkout -b develop
git push -u origin develop

# Создание feature ветки
git checkout -b feature/new-feature
git push -u origin feature/new-feature
```

### 2. Настройка защиты веток

1. Перейдите в **Settings** → **Branches**
2. Добавьте правила для `main` и `develop`:
   - **Require a pull request before merging**
   - **Require status checks to pass before merging**
   - **Require branches to be up to date before merging**
   - **Restrict pushes that create files**

## 📋 **Workflow стратегия**

### 1. Feature Development
```
feature/new-feature → develop (через PR)
```

### 2. Staging Deployment
```
develop → staging (автоматически)
```

### 3. Production Deployment
```
develop → main (через PR) → production (автоматически)
```

## 🚀 **Тестирование CI/CD**

### 1. Тест на develop ветке

```bash
# Создайте тестовый коммит
echo "# Test commit" >> README.md
git add README.md
git commit -m "test: trigger staging deployment"
git push origin develop
```

### 2. Проверка workflow

1. Перейдите в **Actions** вкладку
2. Убедитесь, что workflow запустился
3. Проверьте логи каждого шага

### 3. Тест на main ветке

```bash
# Создайте PR из develop в main
# После merge проверьте production deployment
```

## 🔍 **Мониторинг деплоев**

### 1. GitHub Actions Dashboard

- **Actions** → **All workflows** - общий обзор
- **Actions** → **CI/CD Pipeline** - детали конкретного workflow

### 2. Уведомления

Настройте уведомления в **Settings** → **Notifications**:
- **Actions**: Email при завершении workflow
- **Deployments**: Email при деплое

### 3. Slack интеграция (опционально)

```yaml
# Добавьте в workflow для уведомлений в Slack
- name: Notify Slack
  uses: 8398a7/action-slack@v3
  with:
    status: ${{ job.status }}
    channel: '#deployments'
    webhook_url: ${{ secrets.SLACK_WEBHOOK }}
```

## 🛠️ **Дополнительные настройки**

### 1. Автоматическое создание релизов

```yaml
# Добавьте в workflow после успешного деплоя
- name: Create Release
  uses: actions/create-release@v1
  env:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
  with:
    tag_name: v${{ github.run_number }}
    release_name: Release v${{ github.run_number }}
    body: |
      Changes in this Release
      - Automated deployment to production
    draft: false
    prerelease: false
```

### 2. Автоматическое обновление документации

```yaml
# Добавьте для обновления README с информацией о деплое
- name: Update README
  run: |
    echo "Last deployed: $(date)" >> README.md
    git config --local user.email "action@github.com"
    git config --local user.name "GitHub Action"
    git add README.md
    git commit -m "Update deployment info" || exit 0
    git push
```

### 3. Rollback функциональность

```yaml
# Добавьте job для rollback
rollback:
  runs-on: ubuntu-latest
  if: github.event_name == 'workflow_dispatch' && github.event.inputs.action == 'rollback'
  steps:
    - name: Rollback deployment
      run: |
        export KUBECONFIG=kubeconfig
        helm rollback t-ticker ${{ github.event.inputs.revision }}
```

## 🔧 **Troubleshooting**

### Частые проблемы:

#### 1. Docker Hub rate limit
```yaml
# Добавьте в workflow
- name: Set up Docker Buildx
  uses: docker/setup-buildx-action@v3
  with:
    driver-opts: |
      image=moby/buildkit:buildx-stable-1
```

#### 2. Kubernetes connection timeout
```yaml
# Увеличьте timeout в workflow
- name: Deploy to production
  run: |
    export KUBECONFIG=kubeconfig
    helm upgrade --install t-ticker ./helm/t-ticker \
      --timeout=10m \
      --wait
```

#### 3. Resource constraints
```yaml
# Добавьте resource limits в values
resources:
  limits:
    cpu: 500m
    memory: 512Mi
  requests:
    cpu: 250m
    memory: 256Mi
```

## 📚 **Полезные ссылки**

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Helm Documentation](https://helm.sh/docs/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Docker Hub Documentation](https://docs.docker.com/docker-hub/)

