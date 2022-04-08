# Overview

An expiremental kubectl custom plugin for interacting with OLM's newest file-based catalog (FBC) packaging format in dev environments.

## Quickstart

### Pre-requisites

- Clone the repository: `gh repo clone timflannagan/kubectl-magic-catalog-plugin`
- Install the kubectl plugin locally: `make plugin`

### Verifying the installation

After the repository has been cloned and the plugin Makefile target has been run, verify the plugin is behaving properly by running the following:

```console
$ kubectl catalog --help
A kubectl plugin for creating and managing FBC catalogs in dev environments

Usage:
  catalog [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  create      Instantiate a file-based catalog (FBC) out of thin air
  delete      Delete an existing FBC magic catalog
  help        Help about any command
  update      Update an existing FBC magic catalog

Flags:
      --catalog-name string   Configures the metadata.Name for the generated ConfigMap resource (default "magiccatalog")
  -h, --help                  help for catalog
      --namespace string      Configures the namespace to find the Bundle underlying resources (default "default")

Use "catalog [command] --help" for more information about a command.
```

Note: the catalog command has three distinct sub-commands for working with local file-based catalog files.

### Installing a local FBC

> Note: The instructions below assume you have a Kubernetes and a local FBC available. You can find a sample example-operator.v1.0.0 FBC in the [samples/ directory](./samples/example-operator.v1.0.0.yaml).

First, create a testing namespace that will contain the deployed FBC catalog:

```bash
kubectl create ns test-fbc
```

And then use the kubectl catalog plugin to instantiate a local FBC into that namespace:

```bash
make bin/catalog
```

```bash
./bin/kubectl-catalog create --namespace test-fbc ./samples/example-operator.v1.0.0.yaml
```

And wait until the create command successfully rolls out the generated resources.

Next, verify the installation has a healthy catalog pod that's responsible for serving the grpc connection for that FBC file:

```console
$ kubectl -n test-fbc get pods
NAME               READY   STATUS    RESTARTS   AGE
magiccatalog-pod   1/1     Running   0          34s
```

#### Verifying Catalog Contents using grpcurl

You can use the [grpcurl](https://github.com/fullstorydev/grpcurl) utility to investigate the catalog contents.

First grab the name of the Service resource that sittings in front of the registry Pod:

```console
$ kubectl -n test-fbc get svc
NAME               TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)     AGE
magiccatalog-svc   ClusterIP   10.96.40.134   <none>        50051/TCP   2m23s
```

And port-forward that Service to your localhost 50051 port:

```bash
kubectl -n test-fbc port-forward svc/magiccatalog-svc &
```

And fire off a grpcurl command to verify that the "packageA" in that sample FBC exists:

```console
$ grpcurl -plaintext -d '{"name":"packageA"}' localhost:50051 api.Registry/GetPackage
Handling connection for 50051
{
  "name": "packageA",
  "channels": [
    {
      "name": "stable",
      "csvName": "example-operator.v0.1.0"
    }
  ],
  "defaultChannelName": "stable"
}
```

#### Verify Catalog Contents through OLM APIs

> Note: Before proceeding with the following commands, ensure that [OLM](https://github.com/operator-framework/operator-lifecycle-manager/) is installed on your cluster.

Create the requisite OperatorGroup and Subscriptions resources in the `test-fbc` namespace:

```console
$ cat <<EOF | kubectl create -f -
apiVersion: operators.coreos.com/v1
kind: OperatorGroup
metadata:
  name: test-fbc
  namespace: test-fbc
spec: {}
---
apiVersion: operators.coreos.com/v1alpha1
kind: Subscription
metadata:
  name: test-fbc
  namespace: test-fbc
spec:
  source: magiccatalog
  sourceNamespace: test-fbc
  channel: stable
  name: packageA
EOF

After creating those OLM resources, wait until the example-operator.v0.1.0 ClusterServiceVersion (CSV) resource has been generated, and that resources is reporting a successful installation state:

```console
$ kubectl -n test-fbc get csv example-operator.v0.1.0
NAME                      DISPLAY            VERSION   REPLACES   PHASE
example-operator.v0.1.0   Example Operator   0.1.0                Succeeded
```
