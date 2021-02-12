# Mimik

Simulate being a service (or many) in a mesh. 

Helpful to test Istio features like traffic routing, tracing, security and more. 

## Usage

The following is a fake service mesh deployed with Mimik instances:

![right-lyrics](./docs/examples/mesh.png)

For example, with the following custom resource, the first service can be created as follows:

```yaml
apiVersion: mimik.veicot.io/v1alpha1
kind: Mimik
metadata:
  name: lyrics-page-v1
spec:
  service: lyrics-page
  version: v1
  endpoints:
    - path: /
      method: GET
      connections:
        - service: lyrics-gateway
          port: 8080
          path: songs
          method: GET
    - path: /songs/1
      method: GET
      connections:
        - service: lyrics-gateway
          port: 8080
          path: songs/1
          method: GET
    - path: /health
      method: GET
      connections: []
```

For the rest of the example, follow [this](./docs/examples/operator/operator.md) link.

## Documentation

* [Internals](./docs/internals.md)
* [Usage](./docs/usage.md)