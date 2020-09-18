# Mimik

Simple application to simulate being a microservice in a chain. 

Helpful to test OpenShift Service Mesh (Istio) features like:

* **Traffic Routing**: Allows to set different service versions (v1, v2) to try different routing configurations
* **Tracing**: Propagates the needed HTTP headers for tracing
* **Circuit Breakers**: Returns a configurable HTTP 503 error in any part of the chain
* **JWT**: Propagates the authentication header through services)

Mimik can be instanced many times with different parameters each time thanks to an included OpenShift template.

## Usage

### Template Parameters

* **APP_NAME**: Application name
* **APP_VERSION**: Application version
* **MIMIK_TYPE**: *passthrough* to continue the chain calling another service or *edge* to end it
* **MIMIK_DESTINATION**: The next service's URL to call to
* **MIMIK_SIMULATE_ERROR**: "true" to return a 503 error

## Demonstration

The goal is to deploy some mimik services to achieve the following topology:

![Mesh](mesh.png)

Create the following resources in an OpenShift cluster:

```bash
oc create namespace musik

oc process -f mimik-template.yaml -n musik \
    -p APP_NAME=page \
    -p APP_VERSION=v1 \
    -p MIMIK_TYPE=passthrough \
    -p MIMIK_DESTINATION=http://albums:5000/albums | oc apply -f - -n musik

oc process -f mimik-template.yaml -n musik \
    -p APP_NAME=albums \
    -p APP_VERSION=v1 \
    -p MIMIK_TYPE=passthrough \
    -p MIMIK_DESTINATION=http://songs:5000/songs | oc apply -f - -n musik

oc process -f mimik-template.yaml -n musik \
    -p APP_NAME=songs \
    -p APP_VERSION=v1 \
    -p MIMIK_TYPE=passthrough \
    -p MIMIK_DESTINATION=http://lyrics:5000/lyrics | oc apply -f - -n musik

oc process -f mimik-template.yaml -n musik \
    -p APP_NAME=songs \
    -p APP_VERSION=v2 \
    -p MIMIK_TYPE=passthrough \
    -p MIMIK_DESTINATION=http://lyrics:5000/lyrics | oc apply -f - -n musik

oc process -f mimik-template.yaml -n musik \
    -p APP_NAME=lyrics \
    -p APP_VERSION=v1 \
    -p MIMIK_TYPE=edge | oc apply -f - -n musik

oc apply -f istio.yaml -n musik
```

Then, test the call to the first application (the page):

```bash
curl $(oc get route istio-ingressgateway -o jsonpath='{.spec.host}' -n istio-system)/page
````
