# Usage

Mimik can be instanced many times to create a mesh, each instance needs some configuration to listen for connections and to communicate with other instances.

## Operator

The easiest way to deploy an instance is with Mimik operator.

In Mimik operator's repository there are instructions for installing it.

Having the operator installed, the following CustomResource creates an instance:

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

## Helm

Another way to deploy an instance (without requiring cluster admin permissions) is with a Helm chart:

```bash
echo '[
    {
        "name": "Hello World",
        "path": "/",
        "method": "GET",
        "connections": []
    }
]' > hello-world-v1.json

kubectl create namespace hello-world
kubectl create configmap hello-world-v1 --from-file=hello-world-v1.json -n hello-world

helm install hello-world-v1 ./chart --set version=v1 --set serviceName=hello-world -n hello-world
```

Note: The ConfigMap needs to be called serviceName-version by convention.