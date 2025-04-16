ROOT_PATH := ./../..
K8S_PATH := $(ROOT_PATH)/arch/k8s
K8S_MONITORING_PATH := $(K8S_PATH)/monitoring
HELM_PATH := $(ROOT_PATH)/arch/helm

PINGER_DOCKER := ./build/Dockerfile
PINGER_K8S_PATH := $(K8S_PATH)/pinger
MONITORING_SERVICE_PATH := $(PINGER_K8S_PATH)/monitoring

MIGRATION_PATH := ./db/migrations

MICROSERVICES_NAMESPACE := microservices
PROMETHEUS_NAME := pinger-prometheus
GRAFANA_NAME := grafana-kb
NODE_PORT := 30010
NODE_IP := $(shell kubectl get nodes --namespace $(MICROSERVICES_NAMESPACE) -o jsonpath="{.items[0].status.addresses[0].address}")
NODE_PORT_PROMETHEUS := $(shell kubectl get --namespace $(MICROSERVICES_NAMESPACE) -o jsonpath="{.spec.ports[0].nodePort}" services $(PROMETHEUS_NAME)-kube-pro-prometheus)
NODE_PORT_GRAFANA := $(shell kubectl get --namespace $(MICROSERVICES_NAMESPACE) -o jsonpath="{.spec.ports[0].nodePort}" services $(GRAFANA_NAME))

GRAFANA_USER := admin
GRAFANA_PASS := $(shell kubectl get secret $(GRAFANA_NAME)-admin --namespace $(MICROSERVICES_NAMESPACE) -o jsonpath="{.data.GF_SECURITY_ADMIN_PASSWORD}" | base64 -d)

###### DOCKER ######
build-docker:
	#Build сервиса
	docker build -t vladimirkostin/pinger:latest -f $(PINGER_DOCKER) .;

push-docker:
	#Залить изменения сервиса в репозиторий
	docker push vladimirkostin/pinger;
###### DOCKER ######

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
	kubectl config set-context --current --namespace=$(MICROSERVICES_NAMESPACE)
health:
	curl arch.homework/health/

create_migration:
	./$(ROOT_PATH)/bin/migrate create -ext sql -dir $(MIGRATION_PATH) -seq create_grafana_user

### MONITORING SERVICE ###
add-helm-repo:
	#для настроек мониторинга постгри, надо разобраться в настройке prometheus-postgres-exporter
	-helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
	helm repo update


uninstall-monitoring: uninstall-helms-monitoring delete-config

install-monitoring: install-helms-monitoring install-config

update-grafana: delete-config install-config update-grafana-helm

get-grafana-cred:
	echo $(GRAFANA_USER)
	echo $(GRAFANA_PASS)

get-urls:
	echo "Prometheus URL: http://$(NODE_IP):$(NODE_PORT_PROMETHEUS)/"
	echo "Grafana URL: http://$(NODE_IP):$(NODE_PORT_GRAFANA)/"

get-instances-prometheus:
	#Инстансы prometheus
	kubectl get deployments,pods,services --selector app.kubernetes.io/instance=$(PROMETHEUS_NAME) --namespace $(MICROSERVICES_NAMESPACE)
get-instances-grafana:
	#Инстансы grafana
	kubectl get deployments,pods,services --selector app.kubernetes.io/instance=$(GRAFANA_NAME) --namespace $(MICROSERVICES_NAMESPACE)

update-grafana-helm:
	helm upgrade --namespace $(MICROSERVICES_NAMESPACE) -f $(MONITORING_SERVICE_PATH)/grafana.values.yaml $(GRAFANA_NAME) oci://registry-1.docker.io/bitnamicharts/grafana
delete-config:
	-kubectl delete secret datasources-pinger --namespace $(MICROSERVICES_NAMESPACE)
	-kubectl delete configmap dashboard-kube --namespace $(MICROSERVICES_NAMESPACE)
	-kubectl delete configmap dashboard-postgresql --namespace $(MICROSERVICES_NAMESPACE)
	-kubectl delete configmap dashboard-ingress-nginx --namespace $(MICROSERVICES_NAMESPACE)
	-kubectl delete configmap dashboard-request-handler-performance --namespace $(MICROSERVICES_NAMESPACE)
install-config:
	-kubectl create secret generic datasources-pinger --from-file=$(MONITORING_SERVICE_PATH)/datasources.yaml --namespace $(MICROSERVICES_NAMESPACE)
	-kubectl create configmap dashboard-kube --from-file=$(K8S_MONITORING_PATH)/dashboard-kube.json --namespace $(MICROSERVICES_NAMESPACE)
	-kubectl create configmap dashboard-postgresql --from-file=$(K8S_MONITORING_PATH)/dashboard-postgresql.json --namespace $(MICROSERVICES_NAMESPACE)
	-kubectl create configmap dashboard-ingress-nginx --from-file=$(K8S_MONITORING_PATH)/dashboard-ingress-nginx.json --namespace $(MICROSERVICES_NAMESPACE)
	-kubectl create configmap dashboard-request-handler-performance --from-file=$(K8S_MONITORING_PATH)/dashboard-request-handler-performance.json --namespace $(MICROSERVICES_NAMESPACE)
uninstall-helms-monitoring:
	-helm uninstall $(PROMETHEUS_NAME) --namespace $(MICROSERVICES_NAMESPACE)
	-helm uninstall $(GRAFANA_NAME) --namespace $(MICROSERVICES_NAMESPACE)
install-helms-monitoring:
	-helm install $(PROMETHEUS_NAME) oci://registry-1.docker.io/bitnamicharts/kube-prometheus --namespace $(MICROSERVICES_NAMESPACE) --create-namespace --set prometheus.service.type=NodePort --set prometheus.service.nodePort=$(NODE_PORT)
	-helm install $(GRAFANA_NAME) oci://registry-1.docker.io/bitnamicharts/grafana --values $(MONITORING_SERVICE_PATH)/grafana.values.yaml --namespace $(MICROSERVICES_NAMESPACE)
### MONITORING SERVICE ###

