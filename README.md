# metrics-exporter-playground

## アクセスリンク

- DockerCompose
  - DataSource：http://prometheus:9090
- アクセス
  - ServiceMonitor 確認
    - http://localhost:9090/targets?search=
  - Grafana ダッシュボード
    - http://localhost:3000

## Components

| kube-prometheus-stack |
| :-------------------- |
| Grafana               |
| PushGateway           |
| ServiceMonitor        |
| Prometheus            |
| NodeExporter          |

## 実行

```shell
### 初回時にデプロイする場合
$ kubectl create namespace monitoring
namespace/monitoring created

$ kustomize build ./ --enable-helm | kubectl create -f -
customresourcedefinition.apiextensions.k8s.io/alertmanagerconfigs.monitoring.coreos.com created
customresourcedefinition.apiextensions.k8s.io/alertmanagers.monitoring.coreos.com created
customresourcedefinition.apiextensions.k8s.io/podmonitors.monitoring.coreos.com created
customresourcedefinition.apiextensions.k8s.io/probes.monitoring.coreos.com created
customresourcedefinition.apiextensions.k8s.io/prometheusagents.monitoring.coreos.com created
customresourcedefinition.apiextensions.k8s.io/prometheuses.monitoring.coreos.com created
customresourcedefinition.apiextensions.k8s.io/prometheusrules.monitoring.coreos.com created
customresourcedefinition.apiextensions.k8s.io/scrapeconfigs.monitoring.coreos.com created
customresourcedefinition.apiextensions.k8s.io/servicemonitors.monitoring.coreos.com created
customresourcedefinition.apiextensions.k8s.io/thanosrulers.monitoring.coreos.com created
serviceaccount/kube-prometheus-stack-admission created
serviceaccount/kube-prometheus-stack-alertmanager created
serviceaccount/kube-prometheus-stack-grafana created
serviceaccount/kube-prometheus-stack-kube-state-metrics created
serviceaccount/kube-prometheus-stack-operator created
serviceaccount/kube-prometheus-stack-prometheus created
serviceaccount/kube-prometheus-stack-prometheus-node-exporter created
role.rbac.authorization.k8s.io/kube-prometheus-stack-admission created
role.rbac.authorization.k8s.io/kube-prometheus-stack-grafana created
clusterrole.rbac.authorization.k8s.io/kube-prometheus-stack-admission created
clusterrole.rbac.authorization.k8s.io/kube-prometheus-stack-grafana-clusterrole created
clusterrole.rbac.authorization.k8s.io/kube-prometheus-stack-kube-state-metrics created
clusterrole.rbac.authorization.k8s.io/kube-prometheus-stack-operator created
clusterrole.rbac.authorization.k8s.io/kube-prometheus-stack-prometheus created
rolebinding.rbac.authorization.k8s.io/kube-prometheus-stack-admission created
rolebinding.rbac.authorization.k8s.io/kube-prometheus-stack-grafana created
clusterrolebinding.rbac.authorization.k8s.io/kube-prometheus-stack-admission created
clusterrolebinding.rbac.authorization.k8s.io/kube-prometheus-stack-grafana-clusterrolebinding created
clusterrolebinding.rbac.authorization.k8s.io/kube-prometheus-stack-kube-state-metrics created
clusterrolebinding.rbac.authorization.k8s.io/kube-prometheus-stack-operator created
clusterrolebinding.rbac.authorization.k8s.io/kube-prometheus-stack-prometheus created
configmap/kube-prometheus-stack-alertmanager-overview created
configmap/kube-prometheus-stack-apiserver created
configmap/kube-prometheus-stack-cluster-total created
configmap/kube-prometheus-stack-controller-manager created
configmap/kube-prometheus-stack-etcd created
configmap/kube-prometheus-stack-grafana created
configmap/kube-prometheus-stack-grafana-config-dashboards created
configmap/kube-prometheus-stack-grafana-datasource created
configmap/kube-prometheus-stack-grafana-overview created
configmap/kube-prometheus-stack-k8s-coredns created
configmap/kube-prometheus-stack-k8s-resources-cluster created
configmap/kube-prometheus-stack-k8s-resources-multicluster created
configmap/kube-prometheus-stack-k8s-resources-namespace created
configmap/kube-prometheus-stack-k8s-resources-node created
configmap/kube-prometheus-stack-k8s-resources-pod created
configmap/kube-prometheus-stack-k8s-resources-workload created
configmap/kube-prometheus-stack-k8s-resources-workloads-namespace created
configmap/kube-prometheus-stack-kubelet created
configmap/kube-prometheus-stack-namespace-by-pod created
configmap/kube-prometheus-stack-namespace-by-workload created
configmap/kube-prometheus-stack-node-cluster-rsrc-use created
configmap/kube-prometheus-stack-node-rsrc-use created
configmap/kube-prometheus-stack-nodes created
configmap/kube-prometheus-stack-nodes-darwin created
configmap/kube-prometheus-stack-persistentvolumesusage created
configmap/kube-prometheus-stack-pod-total created
configmap/kube-prometheus-stack-prometheus created
configmap/kube-prometheus-stack-proxy created
configmap/kube-prometheus-stack-scheduler created
configmap/kube-prometheus-stack-workload-total created
secret/alertmanager-kube-prometheus-stack-alertmanager created
secret/kube-prometheus-stack-grafana created
secret/kube-prometheus-stack-prometheus created
service/kube-prometheus-stack-coredns created
service/kube-prometheus-stack-kube-controller-manager created
service/kube-prometheus-stack-kube-etcd created
service/kube-prometheus-stack-kube-proxy created
service/kube-prometheus-stack-kube-scheduler created
service/kube-prometheus-stack-alertmanager created
service/kube-prometheus-stack-grafana created
service/kube-prometheus-stack-kube-state-metrics created
service/kube-prometheus-stack-operator created
service/kube-prometheus-stack-prometheus created
service/kube-prometheus-stack-prometheus-node-exporter created
deployment.apps/kube-prometheus-stack-grafana created
deployment.apps/kube-prometheus-stack-kube-state-metrics created
deployment.apps/kube-prometheus-stack-operator created
daemonset.apps/kube-prometheus-stack-prometheus-node-exporter created
job.batch/kube-prometheus-stack-admission-create created
job.batch/kube-prometheus-stack-admission-patch created
alertmanager.monitoring.coreos.com/kube-prometheus-stack-alertmanager created
prometheus.monitoring.coreos.com/kube-prometheus-stack-prometheus created
prometheusrule.monitoring.coreos.com/kube-prometheus-stack-alertmanager.rules created
prometheusrule.monitoring.coreos.com/kube-prometheus-stack-config-reloaders created
prometheusrule.monitoring.coreos.com/kube-prometheus-stack-etcd created
prometheusrule.monitoring.coreos.com/kube-prometheus-stack-general.rules created
prometheusrule.monitoring.coreos.com/kube-prometheus-stack-k8s.rules.container-cpu-usage-seconds-tot created
prometheusrule.monitoring.coreos.com/kube-prometheus-stack-k8s.rules.container-memory-cache created
prometheusrule.monitoring.coreos.com/kube-prometheus-stack-k8s.rules.container-memory-rss created
prometheusrule.monitoring.coreos.com/kube-prometheus-stack-k8s.rules.container-memory-swap created
prometheusrule.monitoring.coreos.com/kube-prometheus-stack-k8s.rules.container-memory-working-set-by created
prometheusrule.monitoring.coreos.com/kube-prometheus-stack-k8s.rules.container-resource created
prometheusrule.monitoring.coreos.com/kube-prometheus-stack-k8s.rules.pod-owner created
prometheusrule.monitoring.coreos.com/kube-prometheus-stack-kube-apiserver-availability.rules created
prometheusrule.monitoring.coreos.com/kube-prometheus-stack-kube-apiserver-burnrate.rules created
prometheusrule.monitoring.coreos.com/kube-prometheus-stack-kube-apiserver-histogram.rules created
prometheusrule.monitoring.coreos.com/kube-prometheus-stack-kube-apiserver-slos created
prometheusrule.monitoring.coreos.com/kube-prometheus-stack-kube-prometheus-general.rules created
prometheusrule.monitoring.coreos.com/kube-prometheus-stack-kube-prometheus-node-recording.rules created
prometheusrule.monitoring.coreos.com/kube-prometheus-stack-kube-scheduler.rules created
prometheusrule.monitoring.coreos.com/kube-prometheus-stack-kube-state-metrics created
prometheusrule.monitoring.coreos.com/kube-prometheus-stack-kubelet.rules created
prometheusrule.monitoring.coreos.com/kube-prometheus-stack-kubernetes-apps created
prometheusrule.monitoring.coreos.com/kube-prometheus-stack-kubernetes-resources created
prometheusrule.monitoring.coreos.com/kube-prometheus-stack-kubernetes-storage created
prometheusrule.monitoring.coreos.com/kube-prometheus-stack-kubernetes-system created
prometheusrule.monitoring.coreos.com/kube-prometheus-stack-kubernetes-system-apiserver created
prometheusrule.monitoring.coreos.com/kube-prometheus-stack-kubernetes-system-controller-manager created
prometheusrule.monitoring.coreos.com/kube-prometheus-stack-kubernetes-system-kube-proxy created
prometheusrule.monitoring.coreos.com/kube-prometheus-stack-kubernetes-system-kubelet created
prometheusrule.monitoring.coreos.com/kube-prometheus-stack-kubernetes-system-scheduler created
prometheusrule.monitoring.coreos.com/kube-prometheus-stack-node-exporter created
prometheusrule.monitoring.coreos.com/kube-prometheus-stack-node-exporter.rules created
prometheusrule.monitoring.coreos.com/kube-prometheus-stack-node-network created
prometheusrule.monitoring.coreos.com/kube-prometheus-stack-node.rules created
prometheusrule.monitoring.coreos.com/kube-prometheus-stack-prometheus created
prometheusrule.monitoring.coreos.com/kube-prometheus-stack-prometheus-operator created
servicemonitor.monitoring.coreos.com/kube-prometheus-stack-alertmanager created
servicemonitor.monitoring.coreos.com/kube-prometheus-stack-apiserver created
servicemonitor.monitoring.coreos.com/kube-prometheus-stack-coredns created
servicemonitor.monitoring.coreos.com/kube-prometheus-stack-grafana created
servicemonitor.monitoring.coreos.com/kube-prometheus-stack-kube-controller-manager created
servicemonitor.monitoring.coreos.com/kube-prometheus-stack-kube-etcd created
servicemonitor.monitoring.coreos.com/kube-prometheus-stack-kube-proxy created
servicemonitor.monitoring.coreos.com/kube-prometheus-stack-kube-scheduler created
servicemonitor.monitoring.coreos.com/kube-prometheus-stack-kube-state-metrics created
servicemonitor.monitoring.coreos.com/kube-prometheus-stack-kubelet created
servicemonitor.monitoring.coreos.com/kube-prometheus-stack-operator created
servicemonitor.monitoring.coreos.com/kube-prometheus-stack-prometheus created
servicemonitor.monitoring.coreos.com/kube-prometheus-stack-prometheus-node-exporter created
mutatingwebhookconfiguration.admissionregistration.k8s.io/kube-prometheus-stack-admission created
validatingwebhookconfiguration.admissionregistration.k8s.io/kube-prometheus-stack-admission created

$ kustomize build ./ --enable-helm | kubectl create -f -
serviceaccount/prometheus-pushgateway created
service/prometheus-pushgateway created
deployment.apps/prometheus-pushgateway created
```

```shell
### 確認
$ kubectl get pod -n monitoring
NAME                                                        READY   STATUS      RESTARTS   AGE
alertmanager-kube-prometheus-stack-alertmanager-0           2/2     Running     0          21m
kube-prometheus-stack-admission-create-l49gc                0/1     Completed   0          21m
kube-prometheus-stack-admission-patch-7qb9l                 0/1     Completed   0          21m
kube-prometheus-stack-grafana-6f75fcb78f-ck7qb              3/3     Running     0          21m
kube-prometheus-stack-kube-state-metrics-84cfc95b44-62nwn   1/1     Running     0          21m
kube-prometheus-stack-operator-6885c565cb-z7zh4             1/1     Running     0          21m
kube-prometheus-stack-prometheus-node-exporter-xrbmx        1/1     Running     0          21m
prometheus-kube-prometheus-stack-prometheus-0               2/2     Running     0          21m
prometheus-pushgateway-6f9d9b9bf7-92xt4                     1/1     Running     0          25s

$ kubectl get svc -n monitoring
NAME                                             TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)                      AGE
alertmanager-operated                            ClusterIP   None             <none>        9093/TCP,9094/TCP,9094/UDP   21m
kube-prometheus-stack-alertmanager               ClusterIP   10.105.139.95    <none>        9093/TCP,8080/TCP            22m
kube-prometheus-stack-grafana                    ClusterIP   10.96.79.14      <none>        80/TCP                       22m
kube-prometheus-stack-kube-state-metrics         ClusterIP   10.105.59.65     <none>        8080/TCP                     22m
kube-prometheus-stack-operator                   ClusterIP   10.96.31.197     <none>        443/TCP                      22m
kube-prometheus-stack-prometheus                 ClusterIP   10.107.132.79    <none>        9090/TCP,8080/TCP            22m
kube-prometheus-stack-prometheus-node-exporter   ClusterIP   10.109.202.95    <none>        9100/TCP                     22m
prometheus-operated                              ClusterIP   None             <none>        9090/TCP                     21m
prometheus-pushgateway                           ClusterIP   10.104.146.173   <none>        9091/TCP                     38s

### ポートフォワード
$ kubectl port-forward -n monitoring service/kube-prometheus-stack-prometheus 9090:9090
$ kubectl port-forward -n monitoring service/kube-prometheus-stack-grafana 3000:80
$ kubectl port-forward -n sample service/metrics-exporter-sample 8080:8080
$ kubectl port-forward -n monitoring service/prometheus-pushgateway 9091:9091
```

```shell
### 削除
$ kustomize build ./ --enable-helm | kubectl delete -f -
$ kubectl delete namespace monitoring
```

```shell
### トラブルシューティング
$ kubectl get namespace monitoring -o json | jq '.spec.finalizers = []' | kubectl replace --raw /api/v1/namespaces/monitoring/finalize -f -
$ kubectl get namespaces | grep Terminating | awk '{print $1}' | xargs kubectl delete namespace
```

## リソースラベルについて

### Deployment リソース

- `spec.selector.matchLabels.app`
  - ReplicaSet として定義する Pod ラベルを指定
  - 一般に、`spec.template.metadata.labels.app`と同じに設定されることが多い
- `spec.template.metadata.labels.app`
  - Pod として定義するテンプレートラベルを指定
  - Deployment が作成する全ての Pod に適用される
  - 一般に、`spec.selector.matchLabels.app`と同じに設定されることが多い
- `metadata.labels.app`
  - Deployment として公開するリソースに付与するラベルを設定

### Service リソース

- `spec.selectro.app`
  - Service リソースとして公開する ReplicaSet ラベル（Pod ラベル）を指定
  - 例として、Deployment リソースが親となる Pod を定義する場合、Deployment に定義される次の 3 つのラベルは何も同じ値となる
    - `metadata.labels.app`
    - `spec.selector.matchLabels.app`
    - `spec.template.metadata.labels.app`
- `metadata.labels.app`
  - Service として公開するリソースに付与するラベルを設定

### ServiceMonitor リソース

- `spec.selector.matchLabels.app`
  - メトリクスの scrape 先となる Service リソースに設定される `metadata.labels.app` を参照
- `metadata.labels.app`
  - ServiceMonitor として公開するリソースに付与するラベルを設定
  - ServiceMonitor として登録し、Prometheus カスタムリソースにメトリクスを scrape してもらう場合、Prometheus カスタムリソースの `spec.serviceMonitorSelector.matchLabels` に指定された Key-Value を指定

### Prometheus カスタムリソース

- `spec.serviceMonitorSelector.matchLabels`
  - ServiceMonitor として登録する際に使用する Key-Value を定義

## ダッシュボード

<img width="1500" src="https://github.com/GotoRen/metrics-exporter-playground/assets/63791288/038d004a-1fe7-4091-86df-b20dcdb273f4">

<img width="1500" src="https://github.com/GotoRen/metrics-exporter-playground/assets/63791288/64e32ea6-5b43-4497-a279-07d52ddbb59d">

<img width="1500" src="https://github.com/GotoRen/metrics-exporter-playground/assets/63791288/241e1884-9870-4045-b3e4-d0edc12bd757">
