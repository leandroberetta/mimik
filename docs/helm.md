# Instatiate Mimik with Helm

To deploy an instance with a Helm chart:

## Create a ConfigMap with the configuration file for the instance

```bash
kubectl create namespace hello-world

echo '[
    {
        "name": "Hello World",
        "path": "/",
        "method": "GET",
        "connections": []
    }
]' > hello-world-v1.json

kubectl create configmap hello-world-v1 --from-file=hello-world-v1.json -n hello-world
```

Note: The ConfigMap needs to be called serviceName-version by convention.

## Deploy the instance

```bash
helm install hello-world-v1 ./chart --set version=v1 --set serviceName=hello-world -n hello-world
```

