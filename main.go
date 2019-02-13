// Copyright 2019 Layer5.io
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	"github.com/sirupsen/logrus"

	"github.com/layer5io/meshery-istio/istio"
	mesh "github.com/layer5io/meshery-istio/meshes"
)

var (
	gRPCPort = flag.Int("grpc-port", 10000, "The gRPC server port")
)

var log grpclog.LoggerV2

func init() {
	log = grpclog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout)
	grpclog.SetLoggerV2(log)
}

func main() {
	flag.Parse()

	if os.Getenv("DEBUG") == "true" {
		logrus.SetLevel(logrus.DebugLevel)
	}

	addr := fmt.Sprintf(":%d", *gRPCPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logrus.Fatalln("Failed to listen:", err)
	}
	s := grpc.NewServer(
	// grpc.Creds(credentials.NewServerTLSFromCert(&insecure.Cert)),
	)
	mesh.RegisterMeshServiceServer(s, &istio.IstioClient{})

	// Serve gRPC Server
	logrus.Infof("Serving gRPC on %s", addr)
	logrus.Fatal(s.Serve(lis))
}
