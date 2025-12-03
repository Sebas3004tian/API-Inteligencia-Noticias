#!/usr/bin/env bash
set -e

echo "=== Destroy Infra ==="

if ! az account show >/dev/null 2>&1; then
  echo " No estás autenticado. Ejecuta: az login"
  exit 1
fi

SUBSCRIPTION_ID=$(az account show --query id -o tsv)
RESOURCE_GROUP="newsPlatform-rg"
LOCATION="eastus"
STORAGE_NAME="newsplatetfstate"
CONTAINER_NAME="tfstate"
TF_BACKEND_PATH="infra/terraform-backend"
TF_MAIN_PATH="infra"

ACR_NAME="newsplatformacr"
AKS_NAME="newsPlatform-aks"
SP_NAME="ci-cd-sp-newsplatform"

echo " INICIANDO DESTRUCCIÓN DE LA INFRAESTRUCTURA"
echo "---------------------------------------------"


echo " - Eliminando recursos de Kubernetes..."

if az aks show -n "$AKS_NAME" -g "$RESOURCE_GROUP" >/dev/null 2>&1; then
  
  echo " - Obteniendo credenciales..."
  az aks get-credentials -n "$AKS_NAME" -g "$RESOURCE_GROUP" --overwrite-existing || true

  echo " - Eliminando Ingress..."
  kubectl delete -f manifests/k8s/ingress.yaml --ignore-not-found=true || true

  echo " - Eliminando deployments..."
  kubectl delete -f manifests/k8s/deployment-api-news.yaml --ignore-not-found=true || true
  kubectl delete -f manifests/k8s/deployment-embedding-service.yaml --ignore-not-found=true || true
  kubectl delete -f manifests/k8s/deployment-qdrant.yaml --ignore-not-found=true || true

  echo " - Eliminando namespace..."
  kubectl delete -f manifests/k8s/namespace.yaml --ignore-not-found=true || true

  echo " - Eliminando Ingress NGINX (Helm)..."
  helm uninstall ingress-nginx -n ingress-nginx || true
else
  echo " AKS no existe aún. Saltando destrucción de Kubernetes."
fi


echo " - Eliminando imágenes del ACR..."

if az acr show -n "$ACR_NAME" >/dev/null 2>&1; then
  az acr repository delete -n "$ACR_NAME" --repository api-news --yes || true
  az acr repository delete -n "$ACR_NAME" --repository embedding-service --yes || true
else
  echo " ACR no existe. Saltando limpieza."
fi


echo " - Destruyendo infraestructura principal con Terraform..."

if [ -d "$TF_MAIN_PATH" ]; then
  cd "$TF_MAIN_PATH"
  terraform init
  terraform destroy -auto-approve || true
  cd ../
fi


echo " - Destruyendo backend de Terraform..."

if [ -d "$TF_BACKEND_PATH" ]; then
  cd "$TF_BACKEND_PATH"
  terraform init
  terraform destroy -auto-approve || true
  cd ../
fi


echo " - Eliminando Service Principal de CI/CD..."

SP_ID=$(az ad sp list --display-name "$SP_NAME" --query "[0].appId" -o tsv)

if [ -n "$SP_ID" ]; then
  az ad sp delete --id "$SP_ID" || true
  echo " Service Principal eliminado"
else
  echo " No se encontró Service Principal. Nada que eliminar."
fi


echo " Eliminando Resource Group completo..."

if az group show -n "$RESOURCE_GROUP" >/dev/null 2>&1; then
  az group delete -n "$RESOURCE_GROUP" --yes --no-wait
  echo " Resource Group enviado para eliminación"
else
  echo " Resource Group no existe"
fi


echo ""
echo " INFRAESTRUCTURA ELIMINADA COMPLETAMENTE"
echo ""
