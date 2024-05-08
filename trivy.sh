#!/bin/sh

export SERVICE_NAME=${SERVICE}-stage

echo "Pull Docker image to be scanned for vulnerabilities"
echo ${GCR_SERVICE_ACCOUNT} | base64 -d > ./trivy.json
export VERSION=$(head -1 ./semantic_version.txt)

# set image name
export IMAGE_NAME=asia.gcr.io/${GCR_PROJECT_ID}/${SERVICE_NAME}:${VERSION}
# Gcloud auth and check
gcloud auth activate-service-account --key-file ./trivy.json
# Gcloud configuration
gcloud config set project ${GCR_PROJECT_ID}
gcloud config set run/region ${REGION}
# config image registry with gcloud helper
gcloud auth configure-docker asia.gcr.io
# You must be authenticated to Container Registry before tagging an image to successfully push the image to the registry
gcloud auth print-access-token | docker login -u oauth2accesstoken --password-stdin https://asia.gcr.io
echo "pulling docker image..."
# Pulling Docker image
docker pull ${IMAGE_NAME}
# Run vulnerability scanning
trivy image --exit-code 1 ${IMAGE_NAME}
# Gcloud auth revoke
gcloud auth revoke --all
