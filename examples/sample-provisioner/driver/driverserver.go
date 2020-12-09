/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package driver

import (
	"fmt"

	cosi "github.com/kubernetes-sigs/container-object-storage-interface-spec"
	"github.com/minio/minio-go"

	"github.com/minio/minio/pkg/madmin"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	PROVISIONER_NAME = "sample-provisioner.objectstorage.k8s.io"
	VERSION          = "dev"
)

type DriverServer struct {
	Name, Version string
	S3Client      *minio.Client
	S3AdminClient *madmin.AdminClient
}

func (ds *DriverServer) ProvisionerGetInfo(context.Context, *cosi.ProvisionerGetInfoRequest) (*cosi.ProvisionerGetInfoResponse, error) {
	rsp := &cosi.ProvisionerGetInfoResponse{}

	rsp.ProvisionerIdentity = fmt.Sprintf("%s-%s", ds.Name, ds.Version)
	return rsp, nil
}

func (ds DriverServer) ProvisionerCreateBucket(ctx context.Context, req *cosi.ProvisionerCreateBucketRequest) (*cosi.ProvisionerCreateBucketResponse, error) {

	return nil, status.Error(codes.Unavailable, "Method not implemented")
}

func (ds *DriverServer) ProvisionerDeleteBucket(ctx context.Context, req *cosi.ProvisionerDeleteBucketRequest) (*cosi.ProvisionerDeleteBucketResponse, error) {

	return nil, status.Error(codes.Unavailable, "Method not implemented")
}

func (ds *DriverServer) ProvisionerGrantBucketAccess(ctx context.Context, req *cosi.ProvisionerGrantBucketAccessRequest) (*cosi.ProvisionerGrantBucketAccessResponse, error) {
	return nil, status.Error(codes.Unavailable, "Method not implemented")
}

func (ds *DriverServer) ProvisionerRevokeBucketAccess(ctx context.Context, req *cosi.ProvisionerRevokeBucketAccessRequest) (*cosi.ProvisionerRevokeBucketAccessResponse, error) {

	return nil, status.Error(codes.Unavailable, "Method not implemented")
}
