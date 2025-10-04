PROTO_DIR := proto
PROTO_SRC := $(wildcard $(PROTO_DIR)/*.proto)
GO_OUT := .

.PHONY: generate-proto
generate-proto:
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(GO_OUT) \
		--go-grpc_out=$(GO_OUT) \
		$(PROTO_SRC)

#My  My  My
#Development Commands

CLUSTER_NAME=trip-cluster
CLUSTER_REGISTERY_NAME=${CLUSTER_NAME}-registry
API_GATEWAY_IMAGE=api-gateway:1.0.3
TRIP_SERVICE_IMAGE=trip-service:1.0.0
SHELL := /bin/bash

k3d-cluster-up:
	k3d cluster create ${CLUSTER_NAME} --servers 1 --agents 3 -p "80:80@loadbalancer" -p "443:443@loadbalancer" -p "8081:30081@agent:0" -p "8082:30082@agent:1" -p "8083:30083@agent:2" --api-port localhost:6550 \
	--registry-use k3d-${CLUSTER_NAME}-registry
	
k3d-cluster-down:
	k3d cluster delete ${CLUSTER_NAME}

build-web:	
	docker build -t web:1.0.0 -f ./infra/development/docker/web.Dockerfile .

#for path compatibility between normal user build and sudo build "export" added
build-api-gateway-1:
	export PATH=$$PATH:/usr/local/go/bin && \
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -mod=vendor -o ./build/api-gateway ./services/api-gateway  && \
	docker build --no-cache -t ${API_GATEWAY_IMAGE} -f ./infra/development/docker/dev-api-gateway.Dockerfile .

build-api-gateway:
	export PATH=$$PATH:/usr/local/go/bin && \
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -mod=vendor -o ./build/api-gateway ./services/api-gateway  && \
	docker build --no-cache -t ${API_GATEWAY_IMAGE} -f ./infra/development/docker/api-gateway.Dockerfile .

build-trip-service:
	export PATH=$$PATH:/usr/local/go/bin && \
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -mod=vendor -o ./build/trip-service ./services/trip-service/cmd && \
	docker build --no-cache -t ${TRIP_SERVICE_IMAGE} -f ./infra/development/docker/trip-service.Dockerfile .

k3d-create-registry:
	k3d registry create ${CLUSTER_REGISTERY_NAME} --port 5000

k3d-delete-registry:
	k3d registry delete k3d-${CLUSTER_REGISTERY_NAME}

k3d-push-images:
# docker tag web:1.0.0 k3d-${CLUSTER_REGISTERY_NAME}:5000/web:1.0.0 && docker push k3d-${CLUSTER_REGISTERY_NAME}:5000/web:1.0.0
	docker tag web:1.0.0 localhost:5000/web:1.0.0 && docker push localhost:5000/web:1.0.0
	docker tag ${API_GATEWAY_IAMGE} localhost:5000/${API_GATEWAY_IAMGE} && docker push localhost:5000/${API_GATEWAY_IAMGE}

k3d-apply:
	kubectl apply -f ./infra/development/k8s/app-config.yaml
	kubectl apply -f ./infra/development/k8s/api-gateway-deployment.yaml
	kubectl apply -f ./infra/development/k8s/trip-service-deployment.yaml
	kubectl apply -f ./infra/development/k8s/web-deployment.yaml
	kubectl apply -f ./infra/development/k8s/ingress.yaml

go-vendor:
	rm -rf ./vendor 
	go mod vendor

tidy:
	go mod tidy

docker-reset:
	systemctl restart docker