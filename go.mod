module github.com/kubernetes-sigs/container-object-storage-interface-provisioner-sidecar

go 1.14

require (
	github.com/go-ini/ini v1.62.0 // indirect
	github.com/kubernetes-csi/csi-lib-utils v0.9.0
	github.com/kubernetes-sigs/container-object-storage-interface-spec v0.0.0-20201208142312-59e00cb00687
	github.com/minio/minio v0.0.0-20201209163743-e65ed2e44fdd
	github.com/minio/minio-go v6.0.14+incompatible
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
	golang.org/x/net v0.0.0-20200904194848-62affa334b73
	google.golang.org/grpc v1.34.0
	k8s.io/klog v1.0.0
)
