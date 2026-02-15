ifneq ("$(wildcard .env)","")
    include .env
    export # This exports every variable defined in the Makefile to the shell
endif

PROTO_DIR := proto
PROTO_SRC := $(wildcard $(PROTO_DIR)/*.proto)
GO_OUT := .

# pn $$ used for escaping dollar sign - so the makefile does not see it, and shell sees it
# It turns $$ into $ when passing to the shell

.PHONY: generate-proto
generate-proto:
	export PATH=$$PATH:$$(go env GOPATH)/bin && \
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(GO_OUT) \
		--go-grpc_out=$(GO_OUT) \
		$(PROTO_SRC)

clear-proto:
	rm -rf ./shared/proto

#My  My  My
#Development Commands

CLUSTER_NAME=trip-cluster
CLUSTER_REGISTERY_NAME=${CLUSTER_NAME}-registry
#API_GATEWAY_IMAGE=api-gateway:1.0.6
#TRIP_SERVICE_IMAGE=trip-service:1.0.2
#FRONTEND_IMAGE=web:1.0.1
SHELL := /bin/bash

k3d-cluster-up:
	k3d cluster create ${CLUSTER_NAME} --servers 1 --agents 3 -p "8081:8081@loadbalancer" -p "80:80@loadbalancer" -p "443:443@loadbalancer" \
	-p "7071:30071@agent:0" -p "7072:30072@agent:1" -p "7073:30073@agent:2" --api-port localhost:6550 \
	--registry-use k3d-${CLUSTER_NAME}-registry
	
k3d-cluster-down:
	k3d cluster delete ${CLUSTER_NAME}

build-web:	
	docker build -t $(FRONTEND_IMAGE) -f ./infra/development/docker/web.Dockerfile .

#for path compatibility between normal user build and sudo build "export" added
build-api-gateway-dev:
	export PATH=$$PATH:/usr/local/go/bin && \
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -mod=vendor -o ./build/api-gateway ./services/api-gateway  && \
	docker build --no-cache -t ${API_GATEWAY_DEV_IMAGE} -f ./infra/development/docker/api-gateway-dev.Dockerfile .

build-api-gateway:
	export PATH=$$PATH:/usr/local/go/bin && \
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -mod=vendor -o ./build/api-gateway ./services/api-gateway  && \
	docker build --no-cache -t ${API_GATEWAY_IMAGE} -f ./infra/development/docker/api-gateway.Dockerfile .

build-driver-service:
	export PATH=$$PATH:/usr/local/go/bin && \
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -mod=vendor -o ./build/driver-service ./services/driver-service && \
	docker build --no-cache -t ${DRIVER_SERVICE_IMAGE} -f ./infra/development/docker/driver-service.Dockerfile .


build-payment-service:
	export PATH=$$PATH:/usr/local/go/bin && \
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -mod=vendor -o ./build/payment-service ./services/payment-service/cmd && \
	docker build --no-cache -t ${PAYMENT_SERVICE_IMAGE} -f ./infra/development/docker/payment-service.Dockerfile .

build-trip-service:
	export PATH=$$PATH:/usr/local/go/bin && \
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -mod=vendor -o ./build/trip-service ./services/trip-service/cmd && \
	docker build --no-cache -t ${TRIP_SERVICE_IMAGE} -f ./infra/development/docker/trip-service.Dockerfile .

build-all-debug: build-api-gateway-dev
	@echo "All debug builds done!"

#The @ suppresses printing of the command itself, so the output will be just:
build-all: build-trip-service build-api-gateway build-driver-service build-payment-service
	@echo "All builds done!"

debug-api-gateway:
	docker run --rm -d -p 7777:8081 ${API_GATEWAY_IMAGE}

debug-trip-service:
	docker run --rm -d -p 7777:8081 ${TRIP_SERVICE_IMAGE}

rabbitmq:
	docker run --rm -d --name rabbitmq -p 5672:5672 -p 15672:15672 $(RABBITMQ_IMAGE)

env:
	cp ./.env.example ./.env

k3d-create-registry:
	k3d registry create ${CLUSTER_REGISTERY_NAME} --port 5000

k3d-delete-registry:
	k3d registry delete k3d-${CLUSTER_REGISTERY_NAME}

k3d-push-images:
# docker tag web:1.0.0 k3d-${CLUSTER_REGISTERY_NAME}:5000/web:1.0.0 && docker push k3d-${CLUSTER_REGISTERY_NAME}:5000/web:1.0.0
	docker tag $(FRONTEND_IMAGE) localhost:5000/$(FRONTEND_IMAGE) && docker push localhost:5000/$(FRONTEND_IMAGE)
	docker tag $(API_GATEWAY_IMAGE) localhost:5000/$(API_GATEWAY_IMAGE) && docker push localhost:5000/$(API_GATEWAY_IMAGE)

k3d-apply:
	kubectl apply -f ./infra/development/k8s/app-config.yaml
	kubectl apply -f ./infra/development/k8s/api-gateway-deployment.yaml
	kubectl apply -f ./infra/development/k8s/trip-service-deployment.yaml
	kubectl apply -f ./infra/development/k8s/web-deployment.yaml
	kubectl apply -f ./infra/development/k8s/ingress.yaml

vendor:
	rm -rf ./vendor 
	go mod vendor

tidy:
	go mod tidy

docker-reset:
	systemctl restart docker

#helm template <RELEASE_NAME> <CHART_PATH>
helm:
	helm template ${CLUSTER_NAME} ./infra/development/helm/trip-service --output-dir ./build	