# Copyright 2022 IBM Corporation
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

.DEFAULT_GOAL:=help

# Dependence tools
SHELL := /bin/bash
CONTAINER_CLI ?= $(shell basename $(shell which docker))
CONTAINER_BUILD_CMD ?= build
KUBECTL ?= $(shell which kubectl)
OPERATOR_SDK ?= $(shell which operator-sdk)
OPM ?= $(shell which opm)
KUSTOMIZE ?= $(shell which kustomize)
KUSTOMIZE_VERSION=v3.8.7
OPM_VERSION=v1.15.2
YQ_VERSION=3.4.1

# Specify whether this repo is build locally or not, default values is '1';
# If set to 1, then you need to also set 'DOCKER_USERNAME' and 'DOCKER_PASSWORD'
# environment variables before build the repo.
BUILD_LOCALLY ?= 1

VCS_URL = $(shell git config --get remote.origin.url)
VCS_REF ?= $(shell git rev-parse HEAD)
VERSION ?= $(shell cat RELEASE_VERSION)
PREVIOUS_VERSION ?= $(shell cat PREVIOUS_VERSION)
GIT_VERSION ?= $(shell git describe --exact-match 2> /dev/null || \
                 	   git describe --match=$(git rev-parse --short=8 HEAD) --always --dirty --abbrev=8)

# Current Operator image name
OPERATOR_IMAGE_NAME ?= ibm-crossplane-provider-ibm-cloud-operator

# Current Operator bundle image name
BUNDLE_IMAGE_NAME ?= ibm-crossplane-provider-ibm-cloud-operator-bundle

# Options for 'bundle-build'
ifneq ($(origin CHANNELS), undefined)
BUNDLE_CHANNELS := --channels=$(CHANNELS)
endif
ifneq ($(origin DEFAULT_CHANNEL), undefined)
BUNDLE_DEFAULT_CHANNEL := --default-channel=$(DEFAULT_CHANNEL)
endif
BUNDLE_METADATA_OPTS ?= $(BUNDLE_CHANNELS) $(BUNDLE_DEFAULT_CHANNEL)

REGISTRY_DAILY ?= docker-na-public.artifactory.swg-devops.com/hyc-cloud-private-daily-docker-local/ibmcom

ifeq ($(BUILD_LOCALLY),0)
    export CONFIG_DOCKER_TARGET = config-docker
	# Default image repo
	REGISTRY ?= docker-na-public.artifactory.swg-devops.com/hyc-cloud-private-integration-docker-local/ibmcom
else
	REGISTRY ?= docker-na-public.artifactory.swg-devops.com/hyc-cloud-private-scratch-docker-local/ibmcom
endif
OPERATOR_IMAGE := $(REGISTRY)/$(OPERATOR_IMAGE_NAME):$(VERSION)

include common/Makefile.common.mk

############################################################
##@ Develement tools
############################################################

OS    := $(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH  := $(shell uname -m | sed 's/x86_64/amd64/')
OSOPER   := $(shell uname -s | tr '[:upper:]' '[:lower:]' | sed 's/darwin/apple-darwin/' | sed 's/linux/linux-gnu/')
ARCHOPER := $(shell uname -m )

tools: kustomize opm yq ## Install all development tools

kustomize: ## Install kustomize
ifeq (, $(shell which kustomize 2>/dev/null))
	@{ \
	set -e ;\
	mkdir -p bin ;\
	echo "Downloading kustomize ...";\
	curl -sSLo - https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize/$(KUSTOMIZE_VERSION)/kustomize_$(KUSTOMIZE_VERSION)_$(OS)_$(ARCH).tar.gz | tar xzf - -C bin/ ;\
	}
KUSTOMIZE=$(realpath ./bin/kustomize)
else
KUSTOMIZE=$(shell which kustomize)
endif

opm: ## Install operator registry opm
ifeq (, $(shell which opm 2>/dev/null))
	@{ \
	set -e ;\
	mkdir -p bin ;\
	echo "Downloading opm ...";\
	curl -LO https://github.com/operator-framework/operator-registry/releases/download/$(OPM_VERSION)/$(OS)-$(ARCH)-opm ;\
	mv $(OS)-$(ARCH)-opm ./bin/opm ;\
	chmod +x ./bin/opm ;\
	}
OPM=$(realpath ./bin/opm)
else
OPM=$(shell which opm)
endif

yq: ## Install yq, a yaml processor
ifeq (, $(shell which yq 2>/dev/null))
	@{ \
	set -e ;\
	mkdir -p bin ;\
	$(eval ARCH := $(shell uname -m|sed 's/x86_64/amd64/')) \
	echo "Downloading yq ...";\
	curl -LO https://github.com/mikefarah/yq/releases/download/$(YQ_VERSION)/yq_$(OS)_$(ARCH);\
	mv yq_$(OS)_$(ARCH) ./bin/yq ;\
	chmod +x ./bin/yq ;\
	}
YQ=$(realpath ./bin/yq)
else
YQ=$(shell which yq)
endif

ifneq (,$(shell which gsed))
SED=$(shell which gsed)
else
SED=$(shell which sed)
endif


############################################################
##@ Development
############################################################

# artifactory registry
ARTIFACTORY_REGISTRY := docker-na-public.artifactory.swg-devops.com/hyc-cloud-private-scratch-docker-local

ifeq ($(OS),darwin)
	MANIFEST_TOOL_ARGS ?= --username $(DOCKER_USERNAME) --password $(DOCKER_PASSWORD)
else
	MANIFEST_TOOL_ARGS ?= 
endif

check: lint-all ## Check all files lint error
	./common/scripts/lint-csv.sh

install: kustomize ## Install CRDs, controller, and sample CR to a cluster
	$(KUSTOMIZE) build config/development | kubectl apply -f -
	$(KUSTOMIZE) build config/samples | kubectl apply -f -
	- kubectl config set-context --current --namespace=ibm-common-services

uninstall: kustomize ## Uninstall CRDs, controller, and sample CR from a cluster
	$(KUSTOMIZE) build config/samples | kubectl delete --ignore-not-found -f -
	$(KUSTOMIZE) build config/development | kubectl delete --ignore-not-found -f -
	- make clean-cluster

install-catalog-source: ## Install the operator catalog source for testing
	./common/scripts/update_catalogsource.sh $(OPERATOR_IMAGE_NAME) $(REGISTRY)/$(OPERATOR_IMAGE_NAME)-catalog:$(VERSION)

uninstall-catalog-source: ## Uninstall the operator catalog source
	- kubectl -n openshift-marketplace delete catalogsource $(OPERATOR_IMAGE_NAME)

install-operator: install-catalog-source ## Install the operator from catalog source
	- kubectl apply -f config/samples/subscription.yaml 

uninstall-operator: uninstall-catalog-source ## Install the operator from catalog source
	- kubectl get csv -o name | grep ibm-crossplane | xargs kubectl delete
	- kubectl delete --ignore-not-found -f config/samples/subscription.yaml
	- make clean-cluster

clean-cluster: ## Clean up all the resources left in the Kubernetes cluster
	@echo ....... Cleaning up .......
	- kubectl get crossplanes -o name | xargs kubectl delete
	- kubectl get csv -o name | grep ibm-crossplane | xargs kubectl delete
	- kubectl get sub -o name | grep ibm-crossplane | xargs kubectl delete
	- kubectl get installplans | grep ibm-crossplane | awk '{print $$1}' | xargs kubectl delete installplan
	- kubectl get serviceaccounts -o name | grep ibm-crossplane | xargs kubectl delete
	- kubectl get configurationrevisions -o name | xargs kubectl patch -p '{"metadata":{"finalizers": []}}' --type=merge
	- kubectl get configurationrevisions -o name | xargs kubectl delete
	- kubectl get compositeresourcedefinitions -o name | xargs kubectl patch -p '{"metadata":{"finalizers": []}}' --type=merge
	- kubectl get compositeresourcedefinitions -o name | xargs kubectl delete
	- kubectl get configurations -o name --ignore-not-found | xargs kubectl delete
	- kubectl patch locks lock -p '{"metadata":{"finalizers": []}}' --type=merge
	- kubectl get crds,clusterroles,clusterrolebindings -o name | grep crossplane | xargs kubectl delete --ignore-not-found
	- kubectl -n openshift-marketplace get jobs -o name | xargs kubectl -n openshift-marketplace delete --ignore-not-found

create-secret: ## Create artifactory secret in current namespace
	kubectl create secret docker-registry artifactory-secret --docker-server=$(ARTIFACTORY_REGISTRY) --docker-username=$(ARTIFACTORY_USER) --docker-password=$(ARTIFACTORY_TOKEN) --docker-email=none

delete-secret: ## Delete artifactory secret from current namespace
	kubectl delete secret artifactory-secret --ignore-not-found=true

global-pull-secrets: ## Update global pull secrets to use artifactory registries
	./common/scripts/update_global_pull_secrets.sh

############################################################
##@ Build
############################################################

build-controller-image:
	@cd cluster/images/provider-ibm-cloud-controller && make img.build

build-catalog: build-bundle-image build-catalog-source ## Build bundle image and catalog source image for development

# Build bundle image
build-bundle-image: bundle
	@cp -f bundle/manifests/$(OPERATOR_IMAGE_NAME).clusterserviceversion.yaml /tmp/$(OPERATOR_IMAGE_NAME).clusterserviceversion.yaml
	@$(YQ) d -i bundle/manifests/$(OPERATOR_IMAGE_NAME).clusterserviceversion.yaml "spec.replaces"
	@$(SED) -i "s|icr.io/cpopen/cpfs|$(REGISTRY)|g" bundle/manifests/$(OPERATOR_IMAGE_NAME).clusterserviceversion.yaml
	@$(SED) -i "s|icr.io/cpopen|$(REGISTRY)|g" bundle/manifests/$(OPERATOR_IMAGE_NAME).clusterserviceversion.yaml
	$(CONTAINER_CLI) $(CONTAINER_BUILD_CMD) -f bundle.Dockerfile -t $(REGISTRY)/$(BUNDLE_IMAGE_NAME):$(VERSION)-$(ARCH) .
	$(CONTAINER_CLI) push $(REGISTRY)/$(BUNDLE_IMAGE_NAME):$(VERSION)-$(ARCH)
	@mv /tmp/$(OPERATOR_IMAGE_NAME).clusterserviceversion.yaml bundle/manifests/$(OPERATOR_IMAGE_NAME).clusterserviceversion.yaml

# Build catalog source
build-catalog-source: opm build-bundle-image
	$(OPM) -u $(CONTAINER_CLI) index add --bundles $(REGISTRY)/$(BUNDLE_IMAGE_NAME):$(VERSION)-$(ARCH) --tag $(REGISTRY)/$(OPERATOR_IMAGE_NAME)-catalog:$(VERSION)
	$(CONTAINER_CLI) push $(REGISTRY)/$(OPERATOR_IMAGE_NAME)-catalog:$(VERSION)

############################################################
##@ Release
############################################################

copy-operator-data: ## Copy files from submodules before recreating bundle
	cp ./package/crds/* config/crd/bases/

bundle: kustomize copy-operator-data ## Generate bundle manifests and metadata, then validate the generated files
	$(OPERATOR_SDK) generate kustomize manifests -q
	- make bundle-manifests CHANNELS=v3 DEFAULT_CHANNEL=v3

bundle-manifests:
	$(KUSTOMIZE) build config/manifests | $(OPERATOR_SDK) generate bundle \
	-q --overwrite --version $(VERSION) $(BUNDLE_METADATA_OPTS)
	$(OPERATOR_SDK) bundle validate ./bundle
	@./common/scripts/adjust_manifests.sh $(VERSION) $(PREVIOUS_VERSION)


