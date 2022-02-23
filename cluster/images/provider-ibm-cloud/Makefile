# ====================================================================================
# Setup Project

PLATFORMS := linux_amd64 linux_ppc64le linux_s390x
include ../../../build/makelib/common.mk

# ====================================================================================
#  Options
IMAGE_NAME = ibm-crossplane-provider-ibm-cloud
IMAGE = $(BUILD_REGISTRY)/$(IMAGE_NAME):$(VERSION)-$(ARCH)
OSBASEIMAGE = hyc-cloud-private-edge-docker-local.artifactory.swg-devops.com/build-images/ubi8-minimal:latest
include ../../../build/makelib/image.mk

# ====================================================================================
# Targets

img.build:
ifneq (,$(BUILDX))
	$(eval CONTAINER_CLI = docker buildx build --platform linux/$(ARCH) $(BUILDX_ARGS))
else
ifneq (,$(shell which podman))
	$(eval CONTAINER_CLI = podman build --arch $(ARCH) --format docker)
else
ifneq (,$(BUILDX_PLUGIN))
	$(eval CONTAINER_CLI = $(BUILDX_PLUGIN) build --platform linux/$(ARCH) $(BUILDX_ARGS))
else
	$(eval CONTAINER_CLI = docker build)
endif
endif
endif
	@$(INFO) $(CONTAINER_CLI) build $(IMAGE)
	@cp Dockerfile $(IMAGE_TEMP_DIR) || $(FAIL)
	@cp ../../../package/crossplane.yaml $(IMAGE_TEMP_DIR)/package.yaml || $(FAIL)
	@cd $(IMAGE_TEMP_DIR) && $(SED_CMD) 's|BASEIMAGE|$(OSBASEIMAGE)|g' Dockerfile || $(FAIL)
	@cd $(IMAGE_TEMP_DIR) && $(SED_CMD) 's|DOCKER_REGISTRY|$(BUILD_REGISTRY)|g' package.yaml || $(FAIL)
	@cd $(IMAGE_TEMP_DIR) && $(SED_CMD) 's|VERSION|$(VERSION)|g' package.yaml || $(FAIL)
	@$(CONTAINER_CLI) $(BUILD_ARGS) \
		--build-arg ARCH=$(ARCH) \
		--build-arg TINI_VERSION=$(TINI_VERSION) \
		-t $(IMAGE) \
		$(IMAGE_TEMP_DIR) || $(FAIL)
	@$(OK) docker build $(IMAGE)
