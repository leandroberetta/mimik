# Right Lyrics

The following is a fake application represeting a web for search lyrics called Right Lyrics based on Mimik instances.

## Prerequisites 

### Minikube

Start a Minikube instance:

```
minikube start --memory=8g --cpus=4
```

### OLM

Install OLM:

```
minikube addons enable olm
```

### Istio

```
istioctl install --set profile=demo -y

kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.10/samples/addons/prometheus.yaml
kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.10/samples/addons/grafana.yaml
kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.10/samples/addons/jaeger.yaml
kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.10/samples/addons/kiali.yaml
```

### Mimik

Install Mimik operator with OLM:

```
echo "apiVersion: operators.coreos.com/v1alpha1
kind: CatalogSource
metadata:
  name: mimik-catalog
  namespace: olm
spec:
  sourceType: grpc
  image: quay.io/leandroberetta/mimik-operator-index:v0.0.1" | kubectl apply -f -

echo "apiVersion: operators.coreos.com/v1alpha1
kind: Subscription
metadata:
  name: mimik-subscription
  namespace: operators 
spec:
  channel: alpha
  installPlanApproval: Automatic
  name: mimik-operator
  source: mimik-catalog
  sourceNamespace: olm" | kubectl apply -f -
```

## Application

### Deployment

```bash
kubectl create namespace right-lyrics

kubectl apply -f right-lyrics.yaml -n right-lyrics
```

### Topology

![right-lyrics](right-lyrics.png)

### Test

Generate traffic and test the mesh:

```bash
export INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].nodePort}')
export INGRESS_HOST=$(minikube ip)
export GATEWAY_URL=$INGRESS_HOST:$INGRESS_PORT

for i in {1..100}; do curl $GATEWAY_URL/songs/1; done;
```