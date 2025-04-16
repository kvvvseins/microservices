ROOT_PATH := ./../..
K8S_PATH := $(ROOT_PATH)/arch/k8s

PINGER_DOCKER := ./build/Dockerfile
PINGER_K8S_PATH := $(K8S_PATH)/pinger

MIGRATION_PATH := ./db/migrations

MICROSERVICES_NAMESPACE := microservices

###### DOCKER ######
build-docker:
	#Build сервиса
	docker build -t vladimirkostin/pinger:latest -f $(PINGER_DOCKER) .;

push-docker:
	#Залить изменения сервиса в репозиторий
	docker push vladimirkostin/pinger;
###### DOCKER ######

###### INSTALL SERVICE PINGER ######
install: apply migration-up
###### INSTALL SERVICE PINGER ######

###### K8S ######
get-ip-service:
	#В /etc/hosts указать данный IP без порта 192.168.49.2 arch.homework
	minikube service pinger-service -n $(MICROSERVICES_NAMESPACE) --url
	#IP и порт для локального подключения postgresql
	minikube service postgres-pinger -n $(MICROSERVICES_NAMESPACE) --url

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
	kubectl apply -f $(K8S_PATH)/servicemonitor.yaml;
	kubectl apply -f $(K8S_PATH)/ingress.yaml;
	kubectl apply -f $(PINGER_K8S_PATH)/monitoring/pinger-app.yaml

migration-up:
	kubectl apply -f $(PINGER_K8S_PATH)/job/migration-up.yaml;

delete:
	#Удалить манифесты для сервиса
	kubectl delete -f $(K8S_PATH)/ingress.yaml;
	kubectl delete -f $(K8S_PATH)/servicemonitor.yaml;
	kubectl delete -f $(PINGER_K8S_PATH)/service.yaml;
	kubectl delete -f $(PINGER_K8S_PATH)/deployment.yaml;
	kubectl delete -f $(PINGER_K8S_PATH)/postgres.yaml;
	kubectl delete -f $(PINGER_K8S_PATH)/config/app-config.yaml;
	kubectl delete -f $(PINGER_K8S_PATH)/monitoring/pinger-app.yaml
###### K8S ######

setDefaultNamespace:
	kubectl config set-context --current --namespace=$(MICROSERVICES_NAMESPACE)
health:
	curl arch.homework/auth/health/

create_migration:
	./$(ROOT_PATH)/bin/migrate create -ext sql -dir $(MIGRATION_PATH) -seq create_grafana_user


