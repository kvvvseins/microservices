ROOT_PATH := ./../..
K8S_PATH := $(ROOT_PATH)/arch/k8s
HELM_PATH := $(ROOT_PATH)/arch/helm

PINGER_DOCKER := ./build/Dockerfile
PINGER_K8S_PATH := $(K8S_PATH)/pinger

MIGRATION_PATH := ./db/migrations

###### DOCKER ######
build-docker:
	#Build сервиса
	docker build -t vladimirkostin/pinger:latest -f $(PINGER_DOCKER) .;

push-docker:
	#Залить изменения сервиса в репозиторий
	docker push vladimirkostin/pinger;
###### DOCKER ######

###### K8S ######
get-deployment:
	#Удалить манифесты для сервиса
	kubectl get -f $(PINGER_K8S_PATH)/deployment.yaml;
redeploy:
	#Удалить манифесты для сервиса
	kubectl delete -f $(PINGER_K8S_PATH)/deployment.yaml;
	kubectl apply -f $(PINGER_K8S_PATH)/deployment.yaml;

apply:
	#Применить манифесты для сервиса
	kubectl apply -f $(K8S_PATH)/namespace.yaml;
	kubectl apply -f $(PINGER_K8S_PATH)/config/app-config.yaml;
	kubectl apply -f $(PINGER_K8S_PATH)/postgres.yaml;
	kubectl apply -f $(PINGER_K8S_PATH)/deployment.yaml;
	kubectl apply -f $(PINGER_K8S_PATH)/service.yaml;
	kubectl apply -f $(PINGER_K8S_PATH)/ingress.yaml;

migration-up:
	kubectl apply -f $(PINGER_K8S_PATH)/job/migration-up.yaml;

delete:
	#Удалить манифесты для сервиса
	kubectl delete -f $(PINGER_K8S_PATH)/ingress.yaml;
	kubectl delete -f $(PINGER_K8S_PATH)/service.yaml;
	kubectl delete -f $(PINGER_K8S_PATH)/deployment.yaml;
	kubectl delete -f $(PINGER_K8S_PATH)/postgres.yaml;
	kubectl delete -f $(PINGER_K8S_PATH)/config/app-config.yaml;
###### K8S ######

setDefaultNamespace:
	kubectl config set-context --current --namespace=microservices
health:
	curl arch.homework/health/

create_migration:
	./db/migrate create -ext sql -dir $(MIGRATION_PATH) -seq create_users_table

