# ====================================================================================
# Setup Project

PLATFORMS := linux_amd64 linux_ppc64le linux_s390x
include ../../../build/makelib/common.mk

# ====================================================================================
#  Options
IMAGE_NAME = ibm-crossplane-provider-ibm-cloud-controller
IMAGE = $(BUILD_REGISTRY)/$(IMAGE_NAME):$(VERSION)-$(ARCH)
include ../../../build/makelib/image.mk

# ====================================================================================
# Targets

img.build:
ifneq (,$(shell which podman))
	$(eval CONTAINER_CLI = podman build --arch $(ARCH) --format docker)
else
ifneq (,$(BUILDX))
	$(eval CONTAINER_CLI = docker buildx build --platform linux/$(ARCH) $(BUILDX_ARGS))
else
ifneq (,$(BUILDX_PLUGIN))
	$(eval CONTAINER_CLI = $(BUILDX_PLUGIN) build --platform linux/$(ARCH) $(BUILDX_ARGS))
else
	$(eval CONTAINER_CLI = docker build)
endif
endif
endif
	@$(INFO) $(CONTAINER_CLI) $(IMAGE)
	@cp Dockerfile $(IMAGE_TEMP_DIR) || $(FAIL)
	@cp $(OUTPUT_DIR)/bin/$(OS)_$(ARCH)/provider $(IMAGE_TEMP_DIR) || $(FAIL)
	@cd $(IMAGE_TEMP_DIR) && $(SED_CMD) 's|BASEIMAGE|$(OSBASEIMAGE)|g' Dockerfile || $(FAIL)
	@$(CONTAINER_CLI) $(BUILD_ARGS) \
		--build-arg ARCH=$(ARCH) \
		--build-arg TINI_VERSION=$(TINI_VERSION) \
		-t $(IMAGE) \
		$(IMAGE_TEMP_DIR) || $(FAIL)
	@$(OK) $(CONTAINER_CLI) $(IMAGE)
