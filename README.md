# Mimik

Simulate being a service (or many) in a mesh. 

Helpful to test Istio features like traffic routing, tracing, security and more. 

## Example

![right-lyrics](./docs/examples/mesh.png)

## Usage

Having the [Mimik operator](https://github.com/leandroberetta/mimik-operator) installed, the following custom resource creates a Mimik instance:

```yaml
apiVersion: mimik.veicot.io/v1alpha1
kind: Mimik
metadata:
  name: hello-world-v1
spec:
  service: hello-world
  version: v1
  endpoints:
    - path: /hello
      method: GET
      connections: []
```

## Documentation

* [Internals](./docs/internals.md)
* [Usage](./docs/usage.md)
* Examples
    * [Helm](./docs/examples/helm.md)
    * [Operator](./docs/examples/operator.md)
