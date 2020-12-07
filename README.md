# Mimik

Simple application to simulate being a service in a mesh. 

## Usage

Any Mimik instance needs to have the following configuration:

### Environment Variables

The following environment variables are needed to create a Mimik instance:

| Variable | Description |
| - | - |
| MIMIK_SERVICE_NAME | The instance name |
| MIMIK_SERVICE_PORT | The instance port |
| MIMIK_ENDPOINTS_FILE | The file with the endpoints configuration and they connections to upstream services |

### Endpoints

The following file describes the endpoints that a Mimik instances listens to and the connections it has to other upstream services:

```json
[
    {
        "name": "Get songs",
        "path": "/",
        "method": "GET",
        "connections": [
            {
                "name": "songs-service",
                "port": "8080",
                "path": "songs",
                "method": "GET"
            }
        ]
    },
    {
        "name": "Get song with id 1",
        "path": "/songs/1",
        "method": "GET",
        "connections": [
            {
                "name": "songs-service",
                "port": "8080",
                "path": "songs/1",
                "method": "GET"
            },
            {
                "name": "hits-service",
                "port": "8080",
                "path": "hits/1",
                "method": "POST"
            }
        ]
    },
    {
        "name": "Health",
        "path": "/health",
        "method": "GET",
        "connections": []
    }
]
```

## Example

### Right Lyrics

The following command will create a fake application called Right Lyrics, in terms of services it looks like:

![right-lyrics](./example/mesh.png)

#### Deployment

```bash
kubectl create namespace right-lyrics

kubectl label namespace right-lyrics istio-injection=enabled

kubectl create configmap lyrics-page-v1 --from-file=example/lyrics-page-v1.json -n right-lyrics
kubectl create configmap lyrics-gateway-v1 --from-file=example/lyrics-gateway-v1.json -n right-lyrics
kubectl create configmap lyrics-service-v1 --from-file=example/lyrics-service-v1.json -n right-lyrics
kubectl create configmap albums-service-v1 --from-file=example/albums-service-v1.json -n right-lyrics
kubectl create configmap songs-service-v1 --from-file=example/songs-service-v1.json -n right-lyrics
kubectl create configmap songs-service-v2 --from-file=example/songs-service-v2.json -n right-lyrics
kubectl create configmap hits-service-v1 --from-file=example/hits-service-v1.json -n right-lyrics

helm install lyrics-page-v1 ./chart --set version=v1 --set serviceName=lyrics-page -n right-lyrics
helm install lyrics-gateway-v1 ./chart --set version=v1 --set serviceName=lyrics-gateway -n right-lyrics
helm install lyrics-service-v1 ./chart --set version=v1 --set serviceName=lyrics-service -n right-lyrics
helm install albums-service-v1 ./chart --set version=v1 --set serviceName=albums-service -n right-lyrics
helm install songs-service-v1 ./chart --set version=v1 --set serviceName=songs-service -n right-lyrics
helm install songs-service-v2 ./chart --set version=v2 --set serviceName=songs-service --set createService=false -n right-lyrics
helm install hits-service-v1 ./chart --set version=v1 --set serviceName=hits-service -n right-lyrics

export INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].nodePort}')
export INGRESS_HOST=$(minikube ip)
export GATEWAY_URL=$INGRESS_HOST:$INGRESS_PORT 

curl http://$GATEWAY_URL/songs/1
```
