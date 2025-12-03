#!/usr/bin/env bash
set -euo pipefail

echo "=== Build, push and apply Infra ==="

ROOT_DIR=$(cd "$(dirname "$0")/.." && pwd)
ACR_NAME="newsplatformacr"
AKS_NAME="newsPlatform-aks"
RESOURCE_GROUP="newsPlatform-rg"
MANIFESTS_PATH="$ROOT_DIR/manifests/k8s"

echo "== Construyendo imágenes Docker =="

az acr login --name "$ACR_NAME"

docker build -t api-news:latest "$ROOT_DIR/api-news"
docker tag api-news:latest "$ACR_NAME.azurecr.io/api-news:latest"
docker push "$ACR_NAME.azurecr.io/api-news:latest"

docker build -t embedding-service:latest "$ROOT_DIR/embedding-service"
docker tag embedding-service:latest "$ACR_NAME.azurecr.io/embedding-service:latest"
docker push "$ACR_NAME.azurecr.io/embedding-service:latest"


echo "Verificando disponibilidad de imágenes en ACR..."
until az acr repository show --name "$ACR_NAME" --repository api-news >/dev/null 2>&1; do
  echo "Esperando api-news..."
  sleep 3
done
until az acr repository show --name "$ACR_NAME" --repository embedding-service >/dev/null 2>&1; do
  echo "Esperando embedding-service..."
  sleep 3
done

echo "Esperando que nodos de AKS estén listos..."
kubectl wait --for=condition=Ready nodes --timeout=120s || true

echo "Imagenes subidas a la infraestructura correctamente "
