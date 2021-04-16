module sigs.k8s.io/container-object-storage-interface-provisioner-sidecar

go 1.15

require (
	github.com/google/uuid v1.2.0
	github.com/minio/minio v0.0.0-20210415233244-ca9b48b3b423
	github.com/minio/minio-go/v7 v7.0.11-0.20210302210017-6ae69c73ce78
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.0
	golang.org/x/time v0.0.0-20201208040808-7e3f01d25324 // indirect
	google.golang.org/grpc v1.35.0
	k8s.io/api v0.19.4
	k8s.io/apimachinery v0.19.4
	k8s.io/client-go v0.19.4
	k8s.io/klog/v2 v2.2.0
	sigs.k8s.io/container-object-storage-interface-api v0.0.0-20210330175159-2cdabb1a5dc7
	sigs.k8s.io/container-object-storage-interface-spec v0.0.0-20210330184956-b0de747ccee4
)
