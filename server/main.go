/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

//go:generate protoc -I ../helloworld --go_out=plugins=grpc:../helloworld ../helloworld/helloworld.proto

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/naming"
	pb "github.com/surfingPro/gRPCbasis"
	"google.golang.org/grpc"
	"log"
	"net"

	"time"
)

var (
	port = flag.String("port", "50001", "listening port")
)

func main() {
	flag.Parse()

	// 连接etcd集群
	cli, err := clientv3.NewFromURL("http://localhost:2379")
	if err != nil {
		fmt.Println("[测] connect etcd err:", err)
		return
	}

	fmt.Printf("starting hello service at %d", *port)

	// 创建命名解析
	r := &naming.GRPCResolver{Client: cli}
	// 将本服务注册添加etcd中，服务名称为myService，服务地址为本机8001端口
	r.Update(context.TODO(), "myService", naming.Update{Op: naming.Add, Addr: "127.0.0.1:" + string(*port)})

	lis, err := net.Listen("tcp", *port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Printf("%v: Receive is %s\n", time.Now(), in.Name)
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}
