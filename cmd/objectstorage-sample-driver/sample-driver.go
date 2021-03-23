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
	"context"
	"flag"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/minio/minio-go"
	"github.com/minio/minio/pkg/madmin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"k8s.io/klog/v2"

	"sigs.k8s.io/container-object-storage-interface-provisioner-sidecar/pkg/grpcserver"
)

var (
	cosiAddress = "tcp://0.0.0.0:9000"
	endpoint    = "http://0.0.0.0:9000"
	accessKey   = ""
	secretKey   = ""
	ctx         context.Context
)

var cmd = &cobra.Command{
	Use:           os.Args[0],
	Short:         "Sample provisioner for provisioning bucket instance to the backend bucket",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(c *cobra.Command, args []string) error {
		return run(args, cosiAddress)
	},
	DisableFlagsInUseLine: true,
	Version:               VERSION,
}

func init() {
	viper.AutomaticEnv()

	cmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	flag.Set("logtostderr", "true")

	strFlag := func(c *cobra.Command, ptr *string, name string, short string, dfault string, desc string) {
		c.PersistentFlags().
			StringVarP(ptr, name, short, dfault, desc)
	}
	strFlag(cmd, &cosiAddress, "listen-address", "", cosiAddress, "The address for the driver to listen on")
	strFlag(cmd, &endpoint, "s3-endpoint", "", "", "S3-endpont")
	strFlag(cmd, &accessKey, "access-key", "", "", "S3-AccessKey")
	strFlag(cmd, &secretKey, "secret-key", "", "", "S3-SecretKey")
	hideFlag := func(name string) {
		cmd.PersistentFlags().MarkHidden(name)
	}
	hideFlag("alsologtostderr")
	hideFlag("log_backtrace_at")
	hideFlag("log_dir")
	hideFlag("logtostderr")
	hideFlag("master")
	hideFlag("stderrthreshold")
	hideFlag("vmodule")

	// Substitute _ for -
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)

	// suppress the incorrect prefix in glog output
	flag.CommandLine.Parse([]string{})
	viper.BindPFlags(cmd.PersistentFlags())

	var cancel context.CancelFunc

	ctx, cancel = context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSEGV)

	go func() {
		s := <-sigs
		klog.InfoS("Received OS signal", "signal", s.String())
		cancel()
	}()

}

func main() {
	if err := cmd.Execute(); err != nil {
		klog.Fatalln(err)

	}
}

func run(args []string, endpoint string) error {

	u, err := url.Parse(endpoint)
	if err != nil {
		klog.Fatalf("Missing protocol (HTTP/HTTPS) in the endpoint", err.Error())
	}

	secure := u.Scheme == "https"

	klog.V(4).InfoS("Endpoint details", "host", u.Host, "secure", secure)

	minioClient, err := minio.New(u.Host, accessKey, secretKey, secure)
	if err != nil {
		klog.Fatalln(err)
	}

	minioAdminClient, err := madmin.New(u.Host, accessKey, secretKey, secure)
	if err != nil {
		klog.Fatalln(err)
	}

	cds := DriverServer{
		Name:             PROVISIONER_NAME,
		Version:          VERSION,
		MinioClient:      minioClient,
		MinioAdminClient: minioAdminClient,
	}

	ids := IdentityServer{
		Name:    PROVISIONER_NAME,
		Version: VERSION,
	}

	s := grpcserver.NewNonBlockingGRPCServer()
	s.Start(endpoint, &cds, &ids)
	s.Wait()
	return nil
}
