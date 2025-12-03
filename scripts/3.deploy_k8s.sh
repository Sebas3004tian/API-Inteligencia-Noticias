#!/usr/bin/env bash
set -euo pipefail

echo "=== Apply manifests Infra ==="

ROOT_DIR=$(cd "$(dirname "$0")/.." && pwd)
ACR_NAME="newsplatformacr"
AKS_NAME="newsPlatform-aks"
RESOURCE_GROUP="newsPlatform-rg"
MANIFESTS_PATH="$ROOT_DIR/manifests/k8s"

echo "== Configurando credenciales AKS =="
az aks get-credentials -n "$AKS_NAME" -g "$RESOURCE_GROUP" --overwrite-existing

echo "=== Desplegando namespace de la app ==="
kubectl apply -f "$MANIFESTS_PATH/namespace.yaml"

echo "=== Desplegando Servicios y Deployments ==="
kubectl apply -f "$MANIFESTS_PATH/deployment-api-news.yaml"
kubectl apply -f "$MANIFESTS_PATH/deployment-embedding-service.yaml"
kubectl apply -f "$MANIFESTS_PATH/deployment-qdrant.yaml"

echo "=== Instalando ingress-nginx (si no existe) ==="
kubectl get namespace ingress-nginx >/dev/null 2>&1 || kubectl create namespace ingress-nginx

helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update

helm upgrade --install ingress-nginx ingress-nginx/ingress-nginx \
  --namespace ingress-nginx \
  --set controller.replicaCount=1 \
  --set controller.nodeSelector."kubernetes\\.io/os"=linux \
  --set defaultBackend.nodeSelector."kubernetes\\.io/os"=linux

echo "=== Aplicando Ingress ==="
kubectl apply -f "$MANIFESTS_PATH/ingress.yaml"

echo "=== Esperando IP pública del Ingress... ==="
sleep 10

INGRESS_IP=$(kubectl get svc ingress-nginx-controller -n ingress-nginx -o jsonpath="{.status.loadBalancer.ingress[0].ip}")

echo ""
echo "========================================"
echo " Ingress desplegado correctamente"
echo " IP pública: $INGRESS_IP"
echo "========================================"
echo ""
echo "Puedes probar la API en: http://$INGRESS_IP/"

echo ""
echo "=== Generando credenciales requeridas ==="

ACR_USERNAME="$ACR_NAME"

ACR_PASSWORD=$(az acr credential show \
  --name "$ACR_NAME" \
  --query "passwords[0].value" \
  -o tsv)

SUB_ID=$(az account show --query id -o tsv)

AZURE_CREDENTIALS=$(az ad sp create-for-rbac \
  --name "news-sp" \
  --role contributor \
  --scopes "/subscriptions/$SUB_ID/resourceGroups/$RESOURCE_GROUP" \
  --sdk-auth)

echo ""
echo "===================================="
echo " Variables generadas correctamente:"
echo ""
echo "ACR_USERNAME: $ACR_USERNAME"
echo ""
echo "ACR_PASSWORD: $ACR_PASSWORD"
echo ""
echo "AZURE_CREDENTIALS:"
echo "$AZURE_CREDENTIALS"
echo ""
echo "===================================="
