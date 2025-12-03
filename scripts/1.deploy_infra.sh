#!/usr/bin/env bash
set -euo pipefail

echo "=== Deploy de infraestructura Infra ==="

ROOT_DIR=$(cd "$(dirname "$0")/.." && pwd)

TF_BACKEND_PATH="$ROOT_DIR/infra/terraform-backend"
TF_MAIN_PATH="$ROOT_DIR/infra"
MANIFESTS_PATH="$ROOT_DIR/manifests/k8s"

echo "Validando entorno..."

if ! az account show >/dev/null 2>&1; then
  echo "ERROR: No estás autenticado en Azure. Ejecuta: az login"
  exit 1
fi

if [[ ! -d "$TF_BACKEND_PATH" ]]; then
  echo "ERROR: No existe la ruta del backend Terraform: $TF_BACKEND_PATH"
  exit 1
fi

if [[ ! -d "$TF_MAIN_PATH" ]]; then
  echo "ERROR: No existe la ruta de la infraestructura principal: $TF_MAIN_PATH"
  exit 1
fi

if [[ ! -d "$MANIFESTS_PATH" ]]; then
  echo "ERROR: No existe el directorio de manifests Kubernetes: $MANIFESTS_PATH"
  exit 1
fi

echo "Todas las rutas existen. Continuando..."

SUBSCRIPTION_ID=$(az account show --query id -o tsv)
TENANT_ID=$(az account show --query tenantId -o tsv)

RESOURCE_GROUP="newsPlatform-rg"
LOCATION="eastus"

ACR_NAME="newsplatformacr"
AKS_NAME="newsPlatform-aks"
NAMESPACE="news-platform"

SP_NAME="ci-cd-sp-newsplatform"

echo "Subscription: $SUBSCRIPTION_ID"
echo "Tenant:       $TENANT_ID"

echo ""
echo "== Seleccionando suscripción =="
az account set --subscription "$SUBSCRIPTION_ID"


echo "== Verificando Service Principal =="

SP_EXISTE=$(az ad sp list --display-name "$SP_NAME" --query "[0].appId" -o tsv)

if [[ -z "$SP_EXISTE" ]]; then
  echo "Creando SP..."
  SP_OUTPUT=$(az ad sp create-for-rbac \
      --name "$SP_NAME" \
      --role Owner \
      --scopes "/subscriptions/$SUBSCRIPTION_ID")
  
  CLIENT_ID=$(echo "$SP_OUTPUT" | jq -r '.appId')
  CLIENT_SECRET=$(echo "$SP_OUTPUT" | jq -r '.password')

else
  echo "SP ya existe. Regenerando secreto..."
  SP_SECRET=$(az ad sp credential reset --name "$SP_NAME" --query password -o tsv)

  CLIENT_ID="$SP_EXISTE"
  CLIENT_SECRET="$SP_SECRET"
fi

export ARM_SUBSCRIPTION_ID="$SUBSCRIPTION_ID"
export ARM_TENANT_ID="$TENANT_ID"
export ARM_CLIENT_ID="$CLIENT_ID"
export ARM_CLIENT_SECRET="$CLIENT_SECRET"



echo "== Inicializando Backend Terraform =="
cd "$TF_BACKEND_PATH"
terraform init
terraform apply -auto-approve

echo "== Inicializando Infraestructura principal =="
cd "$TF_MAIN_PATH"
terraform init
terraform apply -auto-approve

echo "== Configurando ACR Pull Role =="

echo "Esperando a que el identityProfile del AKS esté disponible..."
for i in {1..20}; do
  KUBELET_ID=$(az aks show -n "$AKS_NAME" -g "$RESOURCE_GROUP" \
      --query identityProfile.kubeletidentity.objectId -o tsv 2>/dev/null || echo "")

  if [[ -n "$KUBELET_ID" && "$KUBELET_ID" != "null" ]]; then
      echo "Kubelet identity lista: $KUBELET_ID"
      break
  fi

  echo "Aún no disponible... reintentando ($i/20)"
  sleep 5
done

if [[ -z "$KUBELET_ID" || "$KUBELET_ID" == "null" ]]; then
  echo "ERROR: identityProfile nunca estuvo disponible."
  exit 1
fi

ACR_ID=$(az acr show --name $ACR_NAME --query id -o tsv)

echo "Asignando rol AcrPull..."
az role assignment create \
  --assignee "$KUBELET_ID" \
  --role "AcrPull" \
  --scope "$ACR_ID" \
  --output none

echo "Esperando propagación de permisos ACR Pull..."
for i in {1..15}; do
  az acr login --name "$ACR_NAME" >/dev/null 2>&1 && break
  echo "Propagando permisos... ($i/15)"
  sleep 4
done

echo "Deploy de infraestructura completado correctamente "
