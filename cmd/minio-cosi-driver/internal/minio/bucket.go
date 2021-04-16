// Copyright 2021 The Kubernetes Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package minio

import (
	"context"
	"github.com/minio/minio/pkg/auth"
	"github.com/minio/minio/pkg/bucket/policy"
	"github.com/minio/minio/pkg/bucket/policy/condition"
	iampolicy "github.com/minio/minio/pkg/iam/policy"
	"k8s.io/klog/v2"

	"github.com/minio/minio-go/v7"
	"github.com/pkg/errors"
)

var ErrBucketAlreadyExists = errors.New("Bucket Already Exists")

type MakeBucketOptions minio.MakeBucketOptions

func (x *C) CreateBucket(ctx context.Context, bucketName string, options MakeBucketOptions) (string, error) {
	if err := x.minioClients.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions(options)); err != nil {
		errCode := minio.ToErrorResponse(err).Code
		if errCode == "BucketAlreadyExists" || errCode == "BucketAlreadyOwnedByYou" {
			return bucketName, ErrBucketAlreadyExists
		}
		return "", err
	}
	return bucketName, nil
}

func (x *C) AddUser(ctx context.Context, bucket string) (*auth.Credentials, error){
	creds, err := auth.GetNewCredentials()
	if err != nil {
		klog.Error("failed to generate new credentails")
		return nil, err
	}

	if err := x.minioClients.adminClient.AddUser(ctx, creds.AccessKey, creds.SecretKey); err != nil {
		klog.Error("failed to create user", err)
		return nil, err
	}

	// Create policy
	p := iampolicy.Policy{
		Version: iampolicy.DefaultVersion,
		Statements: []iampolicy.Statement{
			iampolicy.NewStatement(
				policy.Allow,
				iampolicy.NewActionSet("s3:*"),
				iampolicy.NewResourceSet(iampolicy.NewResource(bucket+"/*", "")),
				condition.NewFunctions(),
			)},
	}

	if err := x.minioClients.adminClient.AddCannedPolicy(context.Background(), "s3:*", &p); err != nil {
		klog.Error("failed to add canned policy", err)
		return nil, err
	}

	if err := x.minioClients.adminClient.SetPolicy(context.Background(), "s3:*", creds.AccessKey, false); err != nil {
		klog.Error("failed to set policy", err)
		return nil, err
	}
	return &creds,nil
}
