# Mimik Operator

Operator to deploy Mimik instances in OpenShift 4 (or any Kubernetes cluster with OLM installed).

## Installation

### OpenShift 4

Create the CatalogSource resource pointing to the Mimik catalog index:

```bash
echo "apiVersion: operators.coreos.com/v1alpha1
kind: CatalogSource
metadata:
  name: mimik-catalog
  namespace: openshift-marketplace
spec:
  sourceType: grpc
  image: quay.io/leandroberetta/mimik-operator-index:v0.0.1" | oc apply -f -
```

Finally create a Subscription resource:

```bash
echo "apiVersion: operators.coreos.com/v1alpha1
kind: Subscription
metadata:
  name: mimik-subscription
  namespace: openshift-operators 
spec:
  channel: alpha
  installPlanApproval: Automatic
  name: mimik-operator
  source: mimik-catalog
  sourceNamespace: openshift-marketplace" | oc apply -f -
```

### Kubernetes

Create the CatalogSource resource pointing to the Mimik catalog index:

```bash
echo "apiVersion: operators.coreos.com/v1alpha1
kind: CatalogSource
metadata:
  name: mimik-catalog
  namespace: olm
spec:
  sourceType: grpc
  image: quay.io/leandroberetta/mimik-operator-index:v0.0.1" | oc apply -f -
```

Finally create a Subscription resource:

```bash
echo "apiVersion: operators.coreos.com/v1alpha1
kind: Subscription
metadata:
  name: mimik-subscription
  namespace: openshift-operators 
spec:
  channel: alpha
  installPlanApproval: Automatic
  name: mimik-operator
  source: mimik-catalog
  sourceNamespace: olm" | oc apply -f -
```

## Usage

For Mimik details and how it works, see the [usage](../docs/usage.md) section.

This is a cluster wide operator so it can create Mimik instances in any namespace with the following CustomResource:

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
      connections:
        - service: hello-world-backend
          port: 8080
          path: hello
          method: GET
```

This is a basic example that creates a Mimik instance that listen for connections at /hello and then tries to connect to another service (upstream connection). You can deploy several instances to create a fake service mesh.

