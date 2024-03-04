# metrics-exporter-playground

## アクセスリンク

- DataSource
  - http://prometheus:9090
- メトリクス確認
  - http://localhost:9090/metrics
- Grafana ダッシュボード
  - http://localhost:3000

## 実行

```shell
### 初回時にデプロイする場合
$ kubectl create namespace monitoring
namespace/monitoring created

$ kustomize build ./ --enable-helm | k create -f -
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
```

```shell
### 確認
$ kubectl get pod
NAME                                                        READY   STATUS      RESTARTS   AGE
alertmanager-kube-prometheus-stack-alertmanager-0           2/2     Running     0          70s
kube-prometheus-stack-admission-create-t5w8k                0/1     Completed   0          71s
kube-prometheus-stack-admission-patch-lgj64                 0/1     Completed   2          71s
kube-prometheus-stack-grafana-6f75fcb78f-rnw45              3/3     Running     0          71s
kube-prometheus-stack-kube-state-metrics-84cfc95b44-gd6kl   1/1     Running     0          71s
kube-prometheus-stack-operator-6885c565cb-nkj6f             1/1     Running     0          71s
kube-prometheus-stack-prometheus-node-exporter-gsf27        1/1     Running     0          71s
prometheus-kube-prometheus-stack-prometheus-0               2/2     Running     0          69s

$ kubectl get svc
NAME                                             TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)                      AGE
alertmanager-operated                            ClusterIP   None            <none>        9093/TCP,9094/TCP,9094/UDP   86s
kube-prometheus-stack-alertmanager               ClusterIP   10.98.184.99    <none>        9093/TCP,8080/TCP            87s
kube-prometheus-stack-grafana                    ClusterIP   10.96.150.251   <none>        80/TCP                       87s
kube-prometheus-stack-kube-state-metrics         ClusterIP   10.100.15.137   <none>        8080/TCP                     87s
kube-prometheus-stack-operator                   ClusterIP   10.110.71.65    <none>        443/TCP                      87s
kube-prometheus-stack-prometheus                 ClusterIP   10.107.15.20    <none>        9090/TCP,8080/TCP            87s
kube-prometheus-stack-prometheus-node-exporter   ClusterIP   10.107.47.153   <none>        9100/TCP                     87s
prometheus-operated                              ClusterIP   None            <none>        9090/TCP                     85s

### ポートフォワード
$ kubectl port-forward -n monitoring service/kube-prometheus-stack-prometheus 9090:9090
$ kubectl port-forward -n monitoring service/kube-prometheus-stack-grafana 3000:80
$ kubectl port-forward -n sample service/metrics-exporter-sample 8080:8080
```

```shell
### 削除
$ kustomize build ./ --enable-helm | k delete -f -
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
  - 一般に、`spec.template.metadata.labels.ap`と同じに設定されることが多い
- `spec.template.metadata.labels.app`
  - Pod として定義するテンプレートラベルを指定
  - Deployment が作成する全ての Pod に適用される
  - 一般に、`spec.selector.matchLabels.app`と同じに設定されることが多い
- `metadata.labels.app`
  - Deployment として公開するリソースに付与するラベルを設定

### Service リソース

- `spec.selectro.app`
  - Service リソースとして公開する ReplicaSet ラベル（Pod ラベル）を指定
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
