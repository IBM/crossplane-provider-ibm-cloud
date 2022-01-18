#!/usr/bin/env bash
#
# Copyright 2021 IBM Corporation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

if [[ (-z "${ARTIFACTORY_REPO}") ]]; then
    ARTIFACTORY_REPO=scratch
fi
if [[ (-z "${ARTIFACTORY_URL}") ]]; then
    ARTIFACTORY_URL="hyc-cloud-private-${ARTIFACTORY_REPO}-docker-local.artifactory.swg-devops.com"
fi
echo "[INFO] ARTIFACTORY_URL set to ${ARTIFACTORY_URL}"

if [[ (-z "${ARTIFACTORY_TOKEN}" || -z "${ARTIFACTORY_USER}") ]]; then
    echo "[ERROR] set env variable ARTIFACTORY_TOKEN and ARTIFACTORY_USER"
    exit 1
fi
echo "[INFO] ARTIFACTORY_TOKEN and ARTIFACTORY_USER set"

if [[ (-z "${IBMCLOUD_API_KEY}") ]]; then
    echo "[ERROR] set env variable IBMCLOUD_API_KEY"
    exit 1
fi
echo "[INFO] IBMCLOUD_API_KEY set to ${IBMCLOUD_API_KEY}"

if [[ (-z "${INSTALL_NAMESPACE}") ]]; then
    INSTALL_NAMESPACE="ibm-common-services"
fi
echo "[INFO] INSTALL_NAMESPACE set to ${INSTALL_NAMESPACE}"

echo "[INFO] create secret with artifactory credentials"
kubectl -n "${INSTALL_NAMESPACE}" create secret docker-registry "artifactory-${ARTIFACTORY_REPO}"\
  --docker-server="${ARTIFACTORY_URL}" \
  --docker-username="${ARTIFACTORY_USER}" \
  --docker-password="${ARTIFACTORY_TOKEN}" \
  --docker-email=none

echo "[INFO] create ibm-crossplane-provider-ibm-cloud ServiceAccount"
PROVIDER_NAME="ibm-crossplane-provider-ibm-cloud"
kubectl -n "${INSTALL_NAMESPACE}" create sa "${PROVIDER_NAME}"

echo "[INFO] apply provider's CRDs"
kubectl apply -f ./config/crd/bases

echo "[INFO] apply provider's RBAC resources"
kubectl -n "${INSTALL_NAMESPACE}" apply \
  -f ./config/rbac/role.yaml \
  -f ./config/rbac/clusterrole.yaml \
  -f ./config/rbac/role_binding.yaml \
  -f ./config/rbac/clusterrole_binding.yaml
yq w ./config/rbac/ibm-crossplane_clusterrole.yaml "metadata.name" "${PROVIDER_NAME}-ibm-crossplane" |\
  kubectl -n "${INSTALL_NAMESPACE}" apply -f -
yq w ./config/rbac/ibm-crossplane_clusterrole_binding.yaml "metadata.name" "${PROVIDER_NAME}-ibm-crossplane" |\
  yq w - "roleRef.name" "${PROVIDER_NAME}-ibm-crossplane" |\
  kubectl -n "${INSTALL_NAMESPACE}" apply -f -

echo "[INFO] create provider's deployment"
sed  "s|quay.io/opencloudio|${ARTIFACTORY_URL}/ibmcom|g" config/manager/manager.yaml |\
  yq w - "metadata.name" "${PROVIDER_NAME}" |\
  yq w - "spec.template.metadata.annotations[olm.targetNamespaces]" "${INSTALL_NAMESPACE}" |\
  kubectl -n "${INSTALL_NAMESPACE}" apply -f -

echo "[INFO] create secret with IBM CLoud API key"
IBMCLOUD_SECRET_NAME="provider-ibm-cloud-secret"
kubectl -n "${INSTALL_NAMESPACE}" create secret generic "${IBMCLOUD_SECRET_NAME}" \
  --from-literal=credentials="${IBMCLOUD_API_KEY}"

echo "[INFO] create ProviderConfig"
cat <<EOF | kubectl apply -f -
apiVersion: ibmcloud.crossplane.io/v1beta1
kind: ProviderConfig
metadata:
  name: ibm-crossplane-provider-ibm-cloud
  namespace: ${INSTALL_NAMESPACE}
spec:
  credentials:
    source: Secret
    secretRef:
      namespace: ${INSTALL_NAMESPACE}
      name: ${IBMCLOUD_SECRET_NAME}
      key: credentials
  region: us-south
EOF
