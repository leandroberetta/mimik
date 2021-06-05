# Install Mimik operator

## OpenShift

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
  image: quay.io/leandroberetta/mimik-operator-index:v0.0.1" | kubectl apply -f -
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
  sourceNamespace: olm" | kubectl apply -f -
```
