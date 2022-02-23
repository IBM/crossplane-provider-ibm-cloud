FROM hyc-cloud-private-edge-docker-local.artifactory.swg-devops.com/build-images/ubi8-minimal:latest

ARG ARCH
ARG TINI_VERSION

ADD provider /usr/local/bin/crossplane-ibm-cloud-provider

EXPOSE 8080
USER 1001
ENTRYPOINT ["crossplane-ibm-cloud-provider"]
