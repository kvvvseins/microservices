ROOT_PATH := ./
K8S_PATH := $(ROOT_PATH)/arch/k8s
HELM_PATH := $(ROOT_PATH)/arch/helm

###### DOCKER ######
login-docker:
	#Логирование в докер
	docker login -u vladimirkostin;
###### DOCKER ######

###### K8S ######
install-nginx-ingress:
	#Установить nginx для ingress
	helm install nginx ingress-nginx/ingress-nginx --namespace ingress-nginx --create-namespace -f $(HELM_PATH)/nginx_ingress.yaml
###### K8S ######
