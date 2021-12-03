# ====================================================================================
# Setup Project

PLATFORMS := linux_amd64 linux_ppc64le linux_s390x
include ../../../build/makelib/common.mk

# ====================================================================================
#  Options
IMAGE_NAME = crossplane-provider-ibm-cloud
IMAGE = $(BUILD_REGISTRY)/$(IMAGE_NAME):$(VERSION)-$(ARCH)
OSBASEIMAGE = scratch
include ../../../build/makelib/image.mk

# ====================================================================================
# Targets

img.build:
ifneq (,$(BUILDX))
	$(eval CONTAINER_CLI = docker buildx build --platform linux/$(ARCH))
else
ifneq (,$(BUILDX_PLUGIN))
	$(eval CONTAINER_CLI = $(BUILDX_PLUGIN) build --platform linux/$(ARCH))
else
	$(eval CONTAINER_CLI = docker build)
endif
endif
	@$(INFO) $(CONTAINER_CLI) build $(IMAGE)
	@cp Dockerfile $(IMAGE_TEMP_DIR) || $(FAIL)
	@cp -R ../../../package $(IMAGE_TEMP_DIR) || $(FAIL)
	@cd $(IMAGE_TEMP_DIR) && $(SED_CMD) 's|BASEIMAGE|$(OSBASEIMAGE)|g' Dockerfile || $(FAIL)
	@cd $(IMAGE_TEMP_DIR) && $(SED_CMD) 's|DOCKER_REGISTRY|$(BUILD_REGISTRY)|g' package/crossplane.yaml || $(FAIL)
	@cd $(IMAGE_TEMP_DIR) && $(SED_CMD) 's|VERSION|$(VERSION)|g' package/crossplane.yaml || $(FAIL)
	@cd $(IMAGE_TEMP_DIR) && find package -type f -name '*.yaml' -exec cat {} >> 'package.yaml' \; -exec printf '\n---\n' \; || $(FAIL)
	@$(CONTAINER_CLI) $(BUILDX_ARGS) $(BUILD_ARGS) \
		--build-arg ARCH=$(ARCH) \
		--build-arg TINI_VERSION=$(TINI_VERSION) \
		-t $(IMAGE) \
		$(IMAGE_TEMP_DIR) || $(FAIL)
	@$(OK) docker build $(IMAGE)
