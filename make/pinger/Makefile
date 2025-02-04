ROOT_PATH := ./../..
K8S_PATH := $(ROOT_PATH)/arch/k8s
HELM_PATH := $(ROOT_PATH)/arch/helm

PINGER_DOCKER := $(ROOT_PATH)/build/pinger/Dockerfile
PINGER_K8S_PATH := $(K8S_PATH)/pinger

###### DOCKER ######
build-docker:
	#Build сервиса pinger
	docker build -t vladimirkostin/pinger:latest -f $(PINGER_DOCKER) $(ROOT_PATH);

push-docker:
	#Залить изменения сервиса pinger в репозиторий
	docker push vladimirkostin/pinger;
###### DOCKER ######

###### K8S ######
apply:
	#Применить манифесты для сервиса pinger
	kubectl apply -f $(K8S_PATH)/namespace.yaml;
	kubectl apply -f $(PINGER_K8S_PATH)/deployment.yaml;
	kubectl apply -f $(PINGER_K8S_PATH)/service.yaml;
	kubectl apply -f $(PINGER_K8S_PATH)/ingress.yaml;

redeploy:
	#Удалить манифесты для сервиса pinger
	kubectl delete -f $(PINGER_K8S_PATH)/deployment.yaml;
	kubectl apply -f $(PINGER_K8S_PATH)/deployment.yaml;

delete:
	#Удалить манифесты для сервиса pinger
	kubectl delete -f $(PINGER_K8S_PATH)/ingress.yaml;
	kubectl delete -f $(PINGER_K8S_PATH)/service.yaml;
	kubectl delete -f $(PINGER_K8S_PATH)/deployment.yaml;
###### K8S ######

health:
	curl arch.homework/health/
