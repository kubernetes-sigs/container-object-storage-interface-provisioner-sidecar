# Deploying Container Object Storage Interface (COSI) Provisioner Sidecar On Kubernetes

This document describes steps for Kubernetes administrators to setup Container Object Storage Interface (COSI) Provisioner Sidecar onto a Kubernetes cluster.

COSI Provisioner Sidecar can be setup using the [kustomization file](https://github.com/kubernetes-retired/cosi-driver-minio/blob/master/kustomization.yaml) from the [cosi-driver-minio](https://github.com/kubernetes-retired/cosi-driver-minio) repository with following command:

```sh
kubectl create -k https://github.com/kubernetes-retired/cosi-driver-minio
```
The output should look like the following:
```sh
namespace/minio-cosi-driver created
serviceaccount/objectstorage-provisioner-sa created
clusterrole.rbac.authorization.k8s.io/objectstorage-provisioner-role created
clusterrolebinding.rbac.authorization.k8s.io/objectstorage-provisioner-role-binding created
configmap/cosi-driver-minio-config created
secret/objectstorage-provisioner created
service/minio created
deployment.apps/minio created
deployment.apps/objectstorage-provisioner created
```

The Provisioner Sidecar will be deployed in the `minio-cosi-driver` namespace.

