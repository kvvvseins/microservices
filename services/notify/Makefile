################### NOTIFY SERVICE ####################

ROOT_PATH := ./../..
K8S_PATH := $(ROOT_PATH)/arch/k8s

SERVICE_DOCKER := ./build/Dockerfile
SERVICE_NAME := notify
SERVICE_K8S_PATH := $(K8S_PATH)/$(SERVICE_NAME)

MIGRATION_PATH := ./db/migrations

MICROSERVICES_NAMESPACE := microservices

###### DOCKER ######

restart-app-with-new-changes: build-docker push-docker redeploy

build-docker:
	#Build сервиса
	docker build -t vladimirkostin/$(SERVICE_NAME):latest -f $(SERVICE_DOCKER) .;

push-docker:
	#Залить изменения сервиса в репозиторий
	docker push vladimirkostin/$(SERVICE_NAME);

#redeploy приложения
redeploy:
	kubectl delete -f $(SERVICE_K8S_PATH)/deployment.yaml;
	kubectl apply -f $(SERVICE_K8S_PATH)/deployment.yaml;
###### DOCKER ######

###### INSTALL SERVICE NOTIFY ######
install: apply
check-status:
	kubectl get pods --namespace $(MICROSERVICES_NAMESPACE) | grep $(SERVICE_NAME)
migration: migration-up
check-ingress:
	kubectl get ingress --namespace $(MICROSERVICES_NAMESPACE)
###### INSTALL SERVICE NOTIFY ######

###### K8S ######
get-ip-service:
	#В /etc/hosts указать данный INTERNAL-IP
	kubectl get nodes -o wide
	#IP и порт для локального подключения postgresql
	minikube service postgres-$(SERVICE_NAME) -n $(MICROSERVICES_NAMESPACE) --url

migration-up:
	kubectl apply -f $(SERVICE_K8S_PATH)/job/migration-up.yaml;

#Применить манифесты для сервиса
apply:
	kubectl apply -f $(SERVICE_K8S_PATH)/config/app-config.yaml;
	kubectl apply -f $(SERVICE_K8S_PATH)/postgres.yaml;
	kubectl apply -f $(SERVICE_K8S_PATH)/deployment.yaml;
	kubectl apply -f $(SERVICE_K8S_PATH)/service.yaml;

#Удалить манифесты для сервиса
delete:
	kubectl delete -f $(SERVICE_K8S_PATH)/service.yaml;
	kubectl delete -f $(SERVICE_K8S_PATH)/deployment.yaml;
	kubectl delete -f $(SERVICE_K8S_PATH)/postgres.yaml;
	kubectl delete -f $(SERVICE_K8S_PATH)/config/app-config.yaml;
###### K8S ######

health:
	curl arch.homework/pc/health/

create_migration:
	./$(ROOT_PATH)/bin/migrate create -ext sql -dir $(MIGRATION_PATH) -seq create_grafana_user
down_migration:
	$(eval DB_URL_FOR_MIGRATION := $(shell minikube service postgres-$(SERVICE_NAME) -n $(MICROSERVICES_NAMESPACE) --url | sed 's/http:\/\///'))
	./$(ROOT_PATH)/bin/migrate -database "postgres://root:passwd@$(DB_URL_FOR_MIGRATION)/$(SERVICE_NAME)?sslmode=disable" -path $(MIGRATION_PATH) down 1


