# Right Lyrics (with Helm)

The following commands will create a fake application called Right Lyrics, in terms of services it looks like:

## Topology

![right-lyrics](../mesh.png)

#### Deployment

```bash
oc create namespace right-lyrics

oc create configmap lyrics-page-v1 --from-file=docs/examples/helm/lyrics-page-v1.json -n right-lyrics
oc create configmap lyrics-gateway-v1 --from-file=docs/examples/helm/lyrics-gateway-v1.json -n right-lyrics
oc create configmap lyrics-service-v1 --from-file=docs/examples/helm/lyrics-service-v1.json -n right-lyrics
oc create configmap albums-service-v1 --from-file=docs/examples/helm/albums-service-v1.json -n right-lyrics
oc create configmap songs-service-v1 --from-file=docs/examples/helm/songs-service-v1.json -n right-lyrics
oc create configmap songs-service-v2 --from-file=docs/examples/helm/songs-service-v2.json -n right-lyrics
oc create configmap hits-service-v1 --from-file=docs/examples/helm/hits-service-v1.json -n right-lyrics

helm install lyrics-page-v1 ./chart --set version=v1 --set serviceName=lyrics-page -n right-lyrics
helm install lyrics-gateway-v1 ./chart --set version=v1 --set serviceName=lyrics-gateway -n right-lyrics
helm install lyrics-service-v1 ./chart --set version=v1 --set serviceName=lyrics-service -n right-lyrics
helm install albums-service-v1 ./chart --set version=v1 --set serviceName=albums-service -n right-lyrics
helm install songs-service-v1 ./chart --set version=v1 --set serviceName=songs-service -n right-lyrics
helm install songs-service-v2 ./chart --set version=v2 --set serviceName=songs-service --set createService=false -n right-lyrics
helm install hits-service-v1 ./chart --set version=v1 --set serviceName=hits-service -n right-lyrics

oc apply -f docs/examples/helm/right-lyrics-gateway.yaml -n right-lyrics
```
