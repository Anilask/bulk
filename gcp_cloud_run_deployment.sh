#!/bin/sh

# Automatic exit from shell script on error
set -e

export ENVIRONMENT=$1
export SERVICE_NAME=${SERVICE}${ENVIRONMENT}

echo "it is a deployment, it should be recorded..."
curl "${CRONITOR_URL}?state=run&series=${SERVICE_NAME}&env=production&message=${SERVICE_NAME}"

echo "Creating Docker image to deploy as Cloud Run service"
echo ${GCR_SERVICE_ACCOUNT} | base64 -d > ./ayopop.json
export VERSION=$(head -1 ./semantic_version.txt)

# set image name
export IMAGE_NAME=asia.gcr.io/${GCR_PROJECT_ID}/${SERVICE_NAME}:${VERSION}
# Gcloud auth and check
gcloud auth activate-service-account --key-file ./ayopop.json
# Gcloud configuration
gcloud config set project ${GCR_PROJECT_ID}
gcloud config set run/region ${REGION}
gcloud config set run/platform managed
# config image registry with gcloud helper
gcloud auth configure-docker asia.gcr.io
# You must be authenticated to Container Registry before tagging an image to successfully push the image to the registry
gcloud auth print-access-token | docker login -u oauth2accesstoken --password-stdin https://asia.gcr.io
echo "creating Docker image..."
# Build Docker image
docker build -t ${IMAGE_NAME} .
# push image to gcr
docker push ${IMAGE_NAME}
# Gcloud auth revoke
gcloud auth revoke --all

echo "Deploying Cloud Run service in GCP"
echo ${GCP_SERVICE_ACCOUNT} | base64 -d > ./ayopop.json

# Gcloud auth and check
gcloud auth activate-service-account --key-file ./ayopop.json

# Gcloud configuration
gcloud config set project ${PROJECT_ID}
gcloud config set compute/region ${REGION}
gcloud config set run/region ${REGION}
gcloud config set run/platform managed

# Ensure the Service Usage API is enabled
is_service_usage_enabled=$(gcloud services list | grep serviceusage.googleapis.com | wc -l)
if [ ${is_service_usage_enabled} == 0 ]; then
    echo "enabling service usage api..."
    gcloud services enable serviceusage.googleapis.com
fi
# Ensure the Cloud Run API is enabled
is_cloud_run_enabled=$(gcloud services list | grep run.googleapis.com | wc -l)
if [ ${is_cloud_run_enabled} == 0 ]; then
    echo "enabling cloud run api..."
    gcloud services enable run.googleapis.com
fi
# Ensure the Serverless VPC Access API is enabled
is_serverless_vpc_enabled=$(gcloud services list | grep vpcaccess.googleapis.com | wc -l)
if [ ${is_serverless_vpc_enabled} == 0 ]; then
    echo "enabling serverless vpc access api..."
    gcloud services enable vpcaccess.googleapis.com
fi
# Ensure the Secret Manager API is enabled
is_secret_api_enabled=$(gcloud services list | grep secretmanager.googleapis.com | wc -l)
if [ ${is_secret_api_enabled} == 0 ]; then
    echo "enabling secret manager api..."
    gcloud services enable secretmanager.googleapis.com
fi
# Ensure the Cloud Resource Manager API is enabled
is_cloud_resource_manager_enabled=$(gcloud services list | grep cloudresourcemanager.googleapis.com | wc -l)
if [ ${is_cloud_resource_manager_enabled} == 0 ]; then
    echo "enabling cloud resource manager api..."
    gcloud services enable cloudresourcemanager.googleapis.com
fi

# Create subnet
is_the_subnet_already_created=$(gcloud compute networks subnets list --filter="(REGION:${REGION})" | grep cloudrunsubnet${ENVIRONMENT} | wc -l)
if [ ${is_the_subnet_already_created} == 0 ]; then
    if [ "${ENVIRONMENT}" == "-dev" ]; then
        echo "define subnet range for the dev environment..."
        export RANGE="10.10.0.0/28"
    fi
    if [ "${ENVIRONMENT}" == "-sandbox" ]; then
        echo "define subnet range for the sandbox environment..."
        export RANGE="10.40.0.0/28"
    fi
    if [ "${ENVIRONMENT}" == "-stage" ]; then
        echo "define subnet range for the stage environment..."
        export RANGE="10.20.0.0/28"
    fi
    if [ "${ENVIRONMENT}" == "" ]; then
        echo "define subnet range the prod environment..."
        export RANGE="10.30.0.0/28"
    fi
    echo "creating subnet..."
    gcloud compute networks subnets create cloudrunsubnet${ENVIRONMENT} --network=${VPC} --region=${REGION} --range=${RANGE}
fi
# Create a Serverless VPC Access connector
is_the_vpc_already_created=$(gcloud compute networks vpc-access connectors list --region ${REGION} | grep serverless-vpc${ENVIRONMENT} | wc -l)
if [ ${is_the_vpc_already_created} == 0 ]; then
    echo "creating vpc access connector for the Cloud Run services..."
    gcloud compute networks vpc-access connectors create serverless-vpc${ENVIRONMENT} --region=${REGION} --subnet=cloudrunsubnet${ENVIRONMENT} --min-instances=2 --max-instances=10 --machine-type=e2-micro
fi
# Configure static outbound IP address for the service
if [ "${ENVIRONMENT}" == "-sandbox" ]; then
    echo "no need for a static outbound IP address in the sandbox environment..."
fi
if [ "${ENVIRONMENT}" != "-sandbox" ]; then
    # Create a Cloud Router
    is_the_cloud_router_already_created=$(gcloud compute routers list --filter="(REGION:${REGION})" | grep cloudrunrouter${ENVIRONMENT} | wc -l)
    if [ ${is_the_cloud_router_already_created} == 0 ]; then
        echo "creating cloud router..."
        # --set-advertisement-mode=custom
        gcloud compute routers create cloudrunrouter${ENVIRONMENT} --project=${PROJECT_ID} --description="Cloud Router for Cloud Run services" --region=${REGION} --network=${VPC}
    fi
    # Create an external IP address
    is_the_ip_already_created=$(gcloud compute addresses list --filter="(REGION:${REGION})" | grep cloudrunstaticoutboundip${ENVIRONMENT} | wc -l)
    if [ ${is_the_ip_already_created} == 0 ]; then
        echo "creating external IP address..."
        gcloud compute addresses create cloudrunstaticoutboundip${ENVIRONMENT} --project=${PROJECT_ID} --description="Outbound static IP address for Cloud Run services" --region=${REGION}
    fi
    # Create a NAT Gateway
    is_the_nat_already_created=$(gcloud compute routers nats list --router=cloudrunrouter${ENVIRONMENT} --region=${REGION} | grep cloudrunnat${ENVIRONMENT} | wc -l)
    if [ ${is_the_nat_already_created} == 0 ]; then
        echo "creating nat gateway..."
        gcloud compute routers nats create cloudrunnat${ENVIRONMENT} --project=${PROJECT_ID} --region=${REGION} --router=cloudrunrouter${ENVIRONMENT} --nat-custom-subnet-ip-ranges=cloudrunsubnet${ENVIRONMENT} --nat-external-ip-pool=cloudrunstaticoutboundip${ENVIRONMENT}
    fi
fi

# Create service account
# NOTE: Service account name must be between 6 and 30 characters. Be careful with very long names
SA_NAME="${SERVICE_NAME:0:27}-sa"
is_the_service_account_already_created=$(gcloud iam service-accounts list --project="${PROJECT_ID}" | grep -E "(^|\s)${SA_NAME}($|\s)" | wc -l)
if [ ${is_the_service_account_already_created} == 0 ]; then
    echo "creating service account..."
    gcloud iam service-accounts create "${SA_NAME}" --display-name="${SA_NAME}"
fi
# Add roles to the service account
IFS=',' read -r -a array <<< "$ROLES"
for ROL in "${array[@]}"
do
    gcloud projects add-iam-policy-binding "${PROJECT_ID}" --member="serviceAccount:${SA_NAME}@${PROJECT_ID}.iam.gserviceaccount.com" --role="$ROL" --user-output-enabled false --condition=None
done

# Decrypting configuration file for the environment
echo ${ANSIBLE_VAULT_PASS} > ./ansible-vault-pass
ansible-vault decrypt ./config/config${ENVIRONMENT}.env --vault-password-file=./ansible-vault-pass
# Create a new secret or add a new versio to an existing secret
secret_already_created=$(gcloud secrets list | grep ${SERVICE}-config${ENVIRONMENT} | wc -l)
if [ ${secret_already_created} == 1 ]; then
    echo "deleting secret..."
    gcloud secrets delete ${SERVICE}-config${ENVIRONMENT} --project=${PROJECT_ID} --quiet
fi
echo "creating new secret..."
gcloud secrets create ${SERVICE}-config${ENVIRONMENT} --data-file=./config/config${ENVIRONMENT}.env

# How to deploy cloud run artifacts
# https://cloud.google.com/sdk/gcloud/reference/beta/run/deploy
# How to access Docker images in other projects
# https://cloud.google.com/run/docs/deploying#other-projects
echo "deploying the Cloud Run service..."
# Check if the service is already deployed
service_already_deployed=$(gcloud beta run services list --platform=managed --region=${REGION} | grep -E "(^|\s)${SERVICE_NAME}($|\s)" | wc -l)
if [ ${service_already_deployed} == 1 ]; then
    echo "service already deployed"
    # Don't serve any traffic when we deploy a new version of the service in the PROD environment
    if [ "${ENVIRONMENT}" == "" ]; then
        echo "The new version won't be serving any traffic"
        export TRAFFIC="--no-traffic"
    else
        echo "The new version will be serving traffic soon"
        export TRAFFIC=""
    fi
else
    echo "service not yet deployed"
    echo "The new version will be serving traffic soon"
    export TRAFFIC=""
fi
# Don't stop the script in case of an error
set +e
gcloud beta run deploy ${SERVICE_NAME} --image ${IMAGE_NAME} ${TRAFFIC} --timeout=${TIMEOUT}s --min-instances=0 --max-instances=${MAX_INSTANCES} --memory=${RAM}Mi --cpu=${CPUS} --concurrency=${CONCURRENCY} --port=${PORT} --service-account=${SA_NAME}@${PROJECT_ID}.iam.gserviceaccount.com --vpc-connector=serverless-vpc${ENVIRONMENT} --ingress=${INGRESS} --vpc-egress=all-traffic --allow-unauthenticated --set-secrets=/app/config/config.env=${SERVICE}-config${ENVIRONMENT}:latest --set-env-vars=OTEL_RESOURCE_ATTRIBUTES=service.name=${SERVICE_NAME}
if [ $? -ne 0 ]; then
    echo "error deploying the Cloud Run service..."
    curl "${CRONITOR_URL}?state=fail&series=${SERVICE_NAME}&env=production&message=${SERVICE_NAME}"
    exit 1
else
    echo "Cloud Run service deployed successfully..."
    curl "${CRONITOR_URL}?state=complete&series=${SERVICE_NAME}&env=production&message=${SERVICE_NAME}"
    exit 0
fi
