SERVICE_VERSION := 1.0.0
SERVICE_NAME := angsuran-service

run:
	go run cmd/main.go

build:
	@echo "Build golang as binary..."
	docker build --no-cache -t ${SERVICE_NAME}:${SERVICE_VERSION} .  
	@echo "Successfully build as binary..."

	@echo "Start saving docker image as .tar file..."
	docker save -o ${SERVICE_NAME}-${SERVICE_VERSION}.tar ${SERVICE_NAME}:${SERVICE_VERSION}
	@echo "Successfully saving image as .tar file..."

	@echo "Start change permission image file to 777..."
	chmod -R 777 ${SERVICE_NAME}-${SERVICE_VERSION}.tar
	@echo "Successfully change permission image file to 777..."

save:
	@echo "Start saving docker image as .tar file..."
	docker save -o ${SERVICE_NAME}-${SERVICE_VERSION}.tar ${SERVICE_NAME}:${SERVICE_VERSION}
	@echo "Successfully saving image as .tar file..."

tidy:
	go mod tidy