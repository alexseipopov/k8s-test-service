# Настройка сервера для деплоя T-Ticker

## 🖥️ **Требования к серверу**

### Минимальные требования:
- **CPU**: 2 ядра
- **RAM**: 4 GB
- **Диск**: 50 GB SSD
- **OS**: Ubuntu 20.04+ / CentOS 8+ / RHEL 8+

### Рекомендуемые требования:
- **CPU**: 4 ядра
- **RAM**: 8 GB
- **Диск**: 100 GB SSD
- **OS**: Ubuntu 22.04 LTS

## 🚀 **Установка Kubernetes кластера**

### Вариант 1: kubeadm (рекомендуется)

#### На всех узлах:

```bash
# Обновление системы
sudo apt update && sudo apt upgrade -y

# Установка Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Отключение swap
sudo swapoff -a
sudo sed -i '/ swap / s/^\(.*\)$/#\1/g' /etc/fstab

# Установка kubeadm, kubelet, kubectl
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
echo "deb https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee /etc/apt/sources.list.d/kubernetes.list
sudo apt update
sudo apt install -y kubelet kubeadm kubectl
sudo apt-mark hold kubelet kubeadm kubectl

# Настройка sysctl
cat <<EOF | sudo tee /etc/modules-load.d/k8s.conf
br_netfilter
EOF

cat <<EOF | sudo tee /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
EOF

sudo sysctl --system
```

#### На master узле:

```bash
# Инициализация кластера
sudo kubeadm init --pod-network-cidr=10.244.0.0/16

# Настройка kubectl для пользователя
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config

# Установка CNI (Flannel)
kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml

# Разрешение запуска подов на master узле (для single-node кластера)
kubectl taint nodes --all node-role.kubernetes.io/control-plane-
```

#### На worker узлах:

```bash
# Присоединение к кластеру (замените на команду из вывода kubeadm init)
sudo kubeadm join <master-ip>:6443 --token <token> --discovery-token-ca-cert-hash <hash>
```

### Вариант 2: MicroK8s (для быстрого тестирования)

```bash
# Установка MicroK8s
sudo snap install microk8s --classic

# Добавление пользователя в группу
sudo usermod -a -G microk8s $USER
newgrp microk8s

# Включение необходимых аддонов
microk8s enable dns storage ingress

# Настройка kubectl
microk8s kubectl config view --raw > ~/.kube/config
```

### Вариант 3: k3s (легковесный Kubernetes)

```bash
# Установка k3s
curl -sfL https://get.k3s.io | sh -

# Настройка kubectl
sudo cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
sudo chown $USER ~/.kube/config
```

## 🔧 **Установка дополнительных компонентов**

### Helm

```bash
# Установка Helm
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
```

### Ingress Controller (NGINX)

```bash
# Установка NGINX Ingress Controller
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm install ingress-nginx ingress-nginx/ingress-nginx \
  --namespace ingress-nginx \
  --create-namespace
```

### Cert-Manager (для SSL сертификатов)

```bash
# Установка cert-manager
helm repo add jetstack https://charts.jetstack.io
helm repo update
helm install cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --version v1.13.0 \
  --set installCRDs=true
```

### Мониторинг (Prometheus + Grafana)

```bash
# Установка Prometheus Stack
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install prometheus prometheus-community/kube-prometheus-stack \
  --namespace monitoring \
  --create-namespace
```

## 🔐 **Настройка безопасности**

### Создание пользователя с ограниченными правами

```bash
# Создание namespace для приложения
kubectl create namespace production

# Создание ServiceAccount
kubectl create serviceaccount t-ticker-deployer -n production

# Создание ClusterRole
cat <<EOF | kubectl apply -f -
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: t-ticker-deployer
rules:
- apiGroups: [""]
  resources: ["pods", "services", "configmaps", "secrets", "persistentvolumeclaims"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
- apiGroups: ["apps"]
  resources: ["deployments", "replicasets", "statefulsets"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
- apiGroups: ["networking.k8s.io"]
  resources: ["ingresses"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
EOF

# Создание ClusterRoleBinding
kubectl create clusterrolebinding t-ticker-deployer \
  --clusterrole=t-ticker-deployer \
  --serviceaccount=production:t-ticker-deployer
```

### Получение kubeconfig для CI/CD

```bash
# Создание токена для ServiceAccount
kubectl create token t-ticker-deployer -n production --duration=8760h

# Получение kubeconfig
kubectl config view --raw --minify > kubeconfig-production.yaml
```

## 🌐 **Настройка DNS и домена**

### Настройка DNS записей

```
# A записи для вашего домена
t-ticker.yourdomain.com     A    <server-ip>
t-ticker-staging.yourdomain.com  A    <server-ip>
```

### Настройка Let's Encrypt

```bash
# Создание ClusterIssuer для Let's Encrypt
cat <<EOF | kubectl apply -f -
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: your-email@example.com
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
    - http01:
        ingress:
          class: nginx
EOF
```

## 📊 **Мониторинг и логирование**

### Установка ELK Stack для логирования

```bash
# Установка Elasticsearch, Logstash, Kibana
helm repo add elastic https://helm.elastic.co
helm install elasticsearch elastic/elasticsearch \
  --namespace logging \
  --create-namespace \
  --set replicas=1 \
  --set resources.requests.memory=1Gi

helm install kibana elastic/kibana \
  --namespace logging \
  --set service.type=NodePort

helm install logstash elastic/logstash \
  --namespace logging
```

### Установка Fluentd для сбора логов

```bash
# Установка Fluentd
helm repo add fluent https://fluent.github.io/helm-charts
helm install fluentd fluent/fluentd \
  --namespace logging \
  --set rbac.create=true
```

## 🔍 **Проверка готовности кластера**

```bash
# Проверка узлов
kubectl get nodes

# Проверка системных подов
kubectl get pods -n kube-system

# Проверка аддонов
kubectl get pods -n ingress-nginx
kubectl get pods -n cert-manager
kubectl get pods -n monitoring

# Тест деплоя
kubectl run test-pod --image=nginx --restart=Never
kubectl get pods
kubectl delete pod test-pod
```

## 📝 **Следующие шаги**

1. **Настройте GitHub Secrets** (см. GITHUB_SETUP.md)
2. **Настройте мониторинг** и алерты
3. **Настройте backup** для PostgreSQL
4. **Настройте логирование** и централизованный сбор логов
5. **Настройте мониторинг безопасности** (Falco, OPA Gatekeeper)

## 🆘 **Устранение неполадок**

### Проблемы с сетью

```bash
# Проверка CNI
kubectl get pods -n kube-system | grep flannel

# Перезапуск CNI
kubectl delete pods -n kube-system -l app=flannel
```

### Проблемы с DNS

```bash
# Проверка DNS
kubectl run test-dns --image=busybox --rm -it --restart=Never -- nslookup kubernetes.default
```

### Проблемы с хранилищем

```bash
# Проверка StorageClass
kubectl get storageclass

# Создание тестового PVC
kubectl apply -f - <<EOF
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: test-pvc
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
EOF
```

