/*
Copyright 2021 The Kubernetes Authors.

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

package main

import (
	"fmt"

	"github.com/minio/minio-go"
	"github.com/minio/minio/pkg/auth"
	"github.com/minio/minio/pkg/bucket/policy"
	"github.com/minio/minio/pkg/bucket/policy/condition"
	iampolicy "github.com/minio/minio/pkg/iam/policy"
	"github.com/minio/minio/pkg/madmin"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"k8s.io/klog/v2"

	cosi "sigs.k8s.io/container-object-storage-interface-spec"
)

type DriverServer struct {
	Name, Version    string
	MinioClient      *minio.Client
	MinioAdminClient *madmin.AdminClient
}

func (ds *DriverServer) ProvisionerGetInfo(context.Context, *cosi.ProvisionerGetInfoRequest) (*cosi.ProvisionerGetInfoResponse, error) {
	rsp := &cosi.ProvisionerGetInfoResponse{}
	rsp.Name = fmt.Sprintf("%s-%s", ds.Name, ds.Version)
	return rsp, nil
}

func (ds DriverServer) ProvisionerCreateBucket(ctx context.Context, req *cosi.ProvisionerCreateBucketRequest) (*cosi.ProvisionerCreateBucketResponse, error) {
	klog.InfoS("Using minio to create bucket")

	if ds.Name == "" {
		return nil, status.Error(codes.Unavailable, ERR_DRIVER_NO_NAME_DEFINED)
	}

	if ds.Version == "" {
		return nil, status.Error(codes.Unavailable, ERR_DRIVER_NO_VERSION_DEFINED)
	}

	s3 := req.Protocol.GetS3()
	if s3 == nil {
		return nil, status.Error(codes.Unavailable, ERR_DRIVER_NO_PROTOCOL_DEFINED)
	}

	err := ds.MinioClient.MakeBucket(s3.BucketName, s3.Region)
	if err != nil {
		// Check to see if the bucket already exists
		exists, errBucketExists := ds.MinioClient.BucketExists(s3.BucketName)
		if errBucketExists == nil && exists {
			klog.InfoS("Bucket already exists", "bucket-name", s3.BucketName)
			return &cosi.ProvisionerCreateBucketResponse{}, nil
		} else {
			klog.ErrorS(err, "Failed to check if bucket already exists", "bucket-name", s3.BucketName)
			return &cosi.ProvisionerCreateBucketResponse{}, err
		}
	}

	klog.InfoS("Successfully created bucket", "bucket-name", s3.BucketName)

	return &cosi.ProvisionerCreateBucketResponse{}, nil
}

func (ds *DriverServer) ProvisionerDeleteBucket(ctx context.Context, req *cosi.ProvisionerDeleteBucketRequest) (*cosi.ProvisionerDeleteBucketResponse, error) {

	s3 := req.Protocol.GetS3()
	if s3 == nil {
		return nil, status.Error(codes.Unavailable, ERR_DRIVER_NO_PROTOCOL_DEFINED)
	}

	klog.InfoS("Deleting bucket", "bucket-name", s3.BucketName)

	if err := ds.MinioClient.RemoveBucket(s3.BucketName); err != nil {
		klog.InfoS("Failed to delete bucket", "bucket-name", s3.BucketName)
		return nil, err
	}

	return &cosi.ProvisionerDeleteBucketResponse{}, nil
}

func (ds *DriverServer) ProvisionerGrantBucketAccess(ctx context.Context, req *cosi.ProvisionerGrantBucketAccessRequest) (*cosi.ProvisionerGrantBucketAccessResponse, error) {
	creds, err := auth.GetNewCredentials()
	if err != nil {
		klog.ErrorS(err, "Failed to generate new credentails")
		return nil, err
	}

	s3 := req.Protocol.GetS3()
	if s3 == nil {
		return nil, status.Error(codes.Unavailable, ERR_DRIVER_NO_PROTOCOL_DEFINED)
	}

	if err := ds.MinioAdminClient.AddUser(context.Background(), creds.AccessKey, creds.SecretKey); err != nil {
		klog.ErrorS(err, "Failed to create user")
		return nil, err
	}

	// Create policy
	p := iampolicy.Policy{
		Version: iampolicy.DefaultVersion,
		Statements: []iampolicy.Statement{
			iampolicy.NewStatement(
				policy.Allow,
				iampolicy.NewActionSet("s3:*"),
				iampolicy.NewResourceSet(iampolicy.NewResource(s3.BucketName+"/*", "")),
				condition.NewFunctions(),
			)},
	}

	if err := ds.MinioAdminClient.AddCannedPolicy(context.Background(), "s3:*", &p); err != nil {
		klog.ErrorS(err, "Failed to add canned policy")
		return nil, err
	}

	if err := ds.MinioAdminClient.SetPolicy(context.Background(), "s3:*", creds.AccessKey, false); err != nil {
		klog.ErrorS(err, "Failed to set policy")
		return nil, err
	}

	return &cosi.ProvisionerGrantBucketAccessResponse{
		Principal:               req.Principal,
		CredentialsFileContents: fmt.Sprintf("[default]\naws_access_key_id %s\naws_secret_access_key %s", creds.AccessKey, creds.SecretKey),
		CredentialsFilePath:     ".aws/credentials",
	}, nil
}

func (ds *DriverServer) ProvisionerRevokeBucketAccess(ctx context.Context, req *cosi.ProvisionerRevokeBucketAccessRequest) (*cosi.ProvisionerRevokeBucketAccessResponse, error) {

	// revokes user access to bucket
	if err := ds.MinioAdminClient.RemoveUser(ctx, req.GetPrincipal()); err != nil {
		klog.ErrorS(err, "Failed to revoke bucket access")
		return nil, err
	}
	return &cosi.ProvisionerRevokeBucketAccessResponse{}, nil
}
