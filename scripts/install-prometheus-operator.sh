helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm create ns monitoring
helm upgrade --install prometheus-operator --namespace=monitoring prometheus-community/prometheus-operator