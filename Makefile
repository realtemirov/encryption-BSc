CURRENT_DIR=$(shell pwd)
APP=$(shell basename ${CURRENT_DIR})
APP_CMD_DIR=${CURRENT_DIR}/cmd
REGISTRY=realtemirov
TAG=latest
ENV_TAG=test
PROJECT_NAME=encryption-bsc
SERVER_USER=SERVER_USER
SERVER_IP=SERVER_HOST

##########################################################################
# Go
go:
	go run cmd/main.go
build:
	CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o ${CURRENT_DIR}/bin/app ${APP_CMD_DIR}/main.go
vendor: 
	go mod vendor
tidy:
	go mod tidy

##########################################################################
# Server
server-copy:
	mkdir -p server
	cp ./Makefile ./server
	cp ./.env ./server
	cp ./Dockerfile ./server
	cp ./docker-compose.yml ./server
	ssh ${SERVER_USER}@${SERVER_IP} "mkdir -p ${PROJECT_NAME}"
	scp -r ./server/.env ${SERVER_USER}@${SERVER_IP}:~/${PROJECT_NAME}/
	scp -r ./server/* ${SERVER_USER}@${SERVER_IP}:~/${PROJECT_NAME}/

##########################################################################
# Image
push-test:
	make server-copy
	sudo docker build -t ${REGISTRY}/${PROJECT_NAME}:${ENV_TAG} .
	sudo docker push ${REGISTRY}/${PROJECT_NAME}:${ENV_TAG} 
	docker rmi ${REGISTRY}/${PROJECT_NAME}:${ENV_TAG}
	ssh ${SERVER_USER}@${SERVER_IP} "mkdir -p ${PROJECT_NAME}"
	scp -r ./server/* ${SERVER_USER}@${SERVER_IP}:~/encryption-bsc
	ssh ${SERVER_USER}@${SERVER_IP} "cd ~/${PROJECT_NAME} && sudo docker compose down && sudo docker compose pull && sudo docker compose up -d"