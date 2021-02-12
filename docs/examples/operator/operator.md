# Right Lyrics (with Mimik operator)

The following commands will create a fake application called Right Lyrics, in terms of services it looks like:

## Topology

![right-lyrics](../mesh.png)

## Deployment

```bash
oc create namespace right-lyrics

oc apply -f right-lyrics.yaml -n right-lyrics
```

Generate traffic and test the mesh:

for i in {1..100}; do curl $(oc get route istio-ingressgateway -o jsonpath='{.spec.host}' -n istio-system)/songs/1; done;