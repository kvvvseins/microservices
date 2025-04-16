ROOT_PATH := .
K8S_PATH := $(ROOT_PATH)/arch/k8s
K8S_MONITORING_PATH := $(K8S_PATH)/monitoring
K8S_MONITORING_DASHBOARD_PATH := $(K8S_MONITORING_PATH)/dashboard
HELM_PATH := $(ROOT_PATH)/arch/helm

MICROSERVICES_NAMESPACE := microservices

###### DOCKER ######
login-docker:
	#Логирование в докер
	docker login -u vladimirkostin;
###### DOCKER ######

###### INSTALL CLUSTER ######
install: pre-install install-monitoring
###### INSTALL CLUSTER ######

###### K8S ######
pre-install: install-nginx-ingress-repo install-nginx-ingress

install-nginx-ingress-repo:
	#Добавляем репу
	-helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx/
	-helm repo update
install-nginx-ingress:
	#Установить nginx для ingress
	-helm install nginx ingress-nginx/ingress-nginx --namespace ingress-nginx --create-namespace -f $(HELM_PATH)/nginx_ingress.yaml
###### K8S ######

### MONITORING SERVICE ###
install-monitoring: add-helm-repo install-config install-helms-monitoring delete-default-dashboard

add-helm-repo:
	#для настроек мониторинга постгри, надо разобраться в настройке prometheus-postgres-exporter
	-helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
	-helm repo update

uninstall-monitoring: uninstall-helms-monitoring delete-config

upgrade-monitoring: upgrade-monitoring-helm delete-default-dashboard

upgrade-monitoring-helm:
	helm upgrade --namespace $(MICROSERVICES_NAMESPACE) prometheus prometheus-community/kube-prometheus-stack -f ./arch/k8s/monitoring/prometheus-grafana.yaml -f ./arch/k8s/monitoring/kube-prometheus.yaml

get-grafana-cred:
	echo admin
	echo $(shell kubectl --namespace $(MICROSERVICES_NAMESPACE) get secrets prometheus-grafana -o jsonpath="{.data.admin-password}" | base64 -d ; echo)

get-urls:
	echo "Prometheus URL: http://prometheus.arch.homework"
	echo "Prometheus URL all metrics: http://prometheus.arch.homework/api/v1/label/__name__/values"
	echo "Grafana URL: http://grafana.arch.homework"

delete-config:
	kubectl delete -f $(K8S_MONITORING_DASHBOARD_PATH)/dashboard-kube.yaml
	kubectl delete -f $(K8S_MONITORING_DASHBOARD_PATH)/dashboard-ingress-nginx.yaml
	kubectl delete -f $(K8S_MONITORING_DASHBOARD_PATH)/dashboard-request-handler-performance.yaml
install-config:
	kubectl apply -f $(K8S_MONITORING_DASHBOARD_PATH)/dashboard-kube.yaml
	kubectl apply -f $(K8S_MONITORING_DASHBOARD_PATH)/dashboard-ingress-nginx.yaml
	kubectl apply -f $(K8S_MONITORING_DASHBOARD_PATH)/dashboard-request-handler-performance.yaml
uninstall-helms-monitoring:
	-helm --namespace $(MICROSERVICES_NAMESPACE) uninstall prometheus
install-helms-monitoring:
	-helm install --namespace $(MICROSERVICES_NAMESPACE) prometheus prometheus-community/kube-prometheus-stack -f ./arch/k8s/monitoring/prometheus-grafana.yaml -f ./arch/k8s/monitoring/kube-prometheus.yaml
delete-default-dashboard:
	-kubectl delete configmap prometheus-kube-prometheus-k8s-coredns
	-kubectl delete configmap prometheus-kube-prometheus-etcd
	-kubectl delete configmap prometheus-kube-prometheus-apiserver
	-kubectl delete configmap prometheus-kube-prometheus-k8s-resources-cluster
	-kubectl delete configmap prometheus-kube-prometheus-k8s-resources-multicluster
	-kubectl delete configmap prometheus-kube-prometheus-k8s-resources-namespace
	-kubectl delete configmap prometheus-kube-prometheus-k8s-resources-node
	-kubectl delete configmap prometheus-kube-prometheus-k8s-resources-pod
	-kubectl delete configmap prometheus-kube-prometheus-k8s-resources-workload
	-kubectl delete configmap prometheus-kube-prometheus-k8s-resources-workloads-namespace
	-kubectl delete configmap prometheus-kube-prometheus-controller-manager
	-kubectl delete configmap prometheus-kube-prometheus-node-cluster-rsrc-use
	-kubectl delete configmap prometheus-kube-prometheus-namespace-by-pod
	-kubectl delete configmap prometheus-kube-prometheus-namespace-by-workload
	-kubectl delete configmap prometheus-kube-prometheus-workload-total
	-kubectl delete configmap prometheus-kube-prometheus-cluster-total
	-kubectl delete configmap prometheus-kube-prometheus-pod-total
	-kubectl delete configmap prometheus-kube-prometheus-persistentvolumesusage
	-kubectl delete configmap prometheus-kube-prometheus-proxy
	-kubectl delete configmap prometheus-kube-prometheus-scheduler
	-kubectl delete configmap prometheus-kube-prometheus-nodes-aix
	-kubectl delete configmap prometheus-kube-prometheus-nodes
	-kubectl delete configmap prometheus-kube-prometheus-prometheus
	-kubectl delete configmap prometheus-kube-prometheus-node-rsrc-use
	-kubectl delete configmap prometheus-kube-prometheus-nodes-darwin
### MONITORING SERVICE ###