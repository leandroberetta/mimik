# Music Store 

The following is a demonstration on how traffic management rules are applied according to the Istio configuration, specifically the virtual services.

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

A very simple microservice architecture will be used to demonstrate the concepts:

### Music Store

This application is generated with Mimik instances:

```
kubectl create namespace music-store
kubectl label namespace music-store istio-injection=enabled
```

#### UI

```
echo "apiVersion: mimik.veicot.io/v1alpha1
kind: Mimik
metadata:
  name: music-store-ui-v1
  namespace: music-store
spec:
  service: music-store-ui
  version: v1
  endpoints:
    - path: /
      method: GET
      connections:
        - service: music-store-backend
          port: 8080
          path: api
          method: GET" | kubectl apply -f -
```

#### Backend

```
echo "apiVersion: mimik.veicot.io/v1alpha1
kind: Mimik
metadata:
  name: music-store-backend-v1
  namespace: music-store
spec:
  service: music-store-backend
  version: v1
  endpoints:
    - path: /api
      method: GET
      connections: []" | kubectl apply -f -
```

## Istio

Configure the following entries in the /etc/hosts to consume the services:

```
echo "$(minikube ip) ui.music.store" | sudo tee -a /etc/hosts
echo "$(minikube ip) backend.music.store" | sudo tee -a /etc/hosts
```

Create the following Istio routing resources:

```
echo "apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: music-store
  namespace: istio-system
spec:
  selector:
    istio: ingressgateway
  servers:
    - port:
        number: 80
        name: http
        protocol: HTTP
      hosts:
        - ui.music.store
        - backend.music.store" | kubectl apply -f -

echo "apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: music-store-ui
  namespace: music-store
spec:
  hosts:
    - ui.music.store
  gateways:
    - istio-system/music-store
  http:
    - match:
        - uri:
            exact: /      
      route:
        - destination:
            host: music-store-ui
            subset: v1
            port:
              number: 8080" | kubectl apply -f -
         
echo "apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: music-store-ui
  namespace: music-store
spec:
  host: music-store-ui
  subsets:
  - name: v1
    labels:
      version: v1" | kubectl apply -f -

echo "apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: music-store-backend
  namespace: music-store
spec:
  hosts:    
    - backend.music.store
  gateways:    
    - istio-system/music-store    
  http:
    - match:
        - uri:
            prefix: /api
      route:
        - destination:
            host: music-store-backend
            subset: v1
            port:
              number: 8080" | kubectl apply -f -
         
echo "apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: music-store-backend
  namespace: music-store
spec:
  host: music-store-backend
  subsets:
  - name: v1
    labels:
      version: v1" | kubectl apply -f -        

```

Generate some load to consume both the UI and the backend from outside the mesh (notice that the UI consumes the backend too):

```
export INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].nodePort}')

fortio load -qps 1 -t 300s http://ui.music.store:$INGRESS_PORT/ &
fortio load -qps 1 -t 300s http://backend.music.store:$INGRESS_PORT/api &
```

Inspect the graph in Kiali and expect the following topology:

![ms1](./ms1.png)

Notice that the backend is consumed by an external client (through the ingress gateway) and also by the UI.

The next step is to add a traffic management rule, for example a fault injection in the backend to start returning a 500 error in every request:

echo "apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: music-store-backend
  namespace: music-store
spec:
  hosts:
    - backend.music.store
  gateways:
    - istio-system/music-store    
  http:
    - fault:
        abort:
          httpStatus: 500
          percentage:
            value: 100
      match:
        - uri:
            prefix: /api
      route:
        - destination:
            host: music-store-backend
            subset: v1
            port:
              number: 8080" | kubectl apply -f -

Inspect the graph in Kiali:

![m2](./ms2.png)

Notice that the rule is applying to the external client only but the internal client (the UI) is working good. 

This behaviour is expected because in the backend's virtual service there are some things missing to apply the rule to internal traffic.

Apply the following virtual service that adds a new entry in the hosts list (the internal service of the backend):

echo "apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: music-store-backend
  namespace: music-store
spec:
  hosts:
    - music-store-backend.music-store.svc.cluster.local
    - backend.music.store
  gateways:    
    - istio-system/music-store    
  http:
    - fault:
        abort:
          httpStatus: 500
          percentage:
            value: 100
      match:
        - uri:
            prefix: /api
      route:
        - destination:
            host: music-store-backend
            subset: v1
            port:
              number: 8080" | kubectl apply -f -              

Inpect the graph in Kiali and observe that the rule is still not applying, it is still working, and that is because another configuration is missing.

In the virtual service, the only gateway that is configured is the gateway that is related to the ingress gateway (external traffic getting into the mesh), so internal traffic is not being controlled by this rule, to fix this situation, an special value "mesh" can be configured in the gateways list as follows:

echo "apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: music-store-backend
  namespace: music-store
spec:
  hosts:
    - music-store-backend.music-store.svc.cluster.local
    - backend.music.store
  gateways:    
    - mesh
    - istio-system/music-store    
  http:
    - fault:
        abort:
          httpStatus: 500
          percentage:
            value: 100
      match:
        - uri:
            prefix: /api
      route:
        - destination:
            host: music-store-backend
            subset: v1
            port:
              number: 8080" | kubectl apply -f -  

Inpect the graph in Kiali and observe that the rule is applying for both external and internal calls:

![ms3](./ms3.png)
