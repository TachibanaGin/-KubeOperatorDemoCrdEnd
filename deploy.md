#0.开始之前
将namespace设置为istio的自动注入
```
kubectl label namespace beijing istio-injection=enabled
```
查看
```
kubectl get namespace -L istio-injection
```
#1.build
```
GOOS=linux go build -o ./app .
```
#2.build-image
```
docker build -t registry.cn-qingdao.aliyuncs.com/fuck-k8s/lhcz-demo-crd-end:v7 . 
```
#3.rolebinding
```
kubectl create clusterrolebinding default-view --clusterrole=view --serviceaccount=beijing:default

```
注:
```
kubectl create clusterrolebinding myapp-view-binding --clusterrole=view --serviceaccount=acme:myapp
```
*在整个集群范围内将view ClusterRole授予名字空间”acme”内的服务账户”myapp”*
#4.deploy
```
kubectl apply -f lhcz-demo-crd.yaml
```