package main

import (
	pb "../helloworld"
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	etcdnaming "github.com/coreos/etcd/clientv3/naming"
	"google.golang.org/grpc"
	"strconv"
	"time"
)

func main() {
	// 连接etcd集群
	cli, err := clientv3.New(clientv3.Config{
		// etcd集群成员节点列表
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("[测] connect etcd err:", err)
		return
	}

	r := &etcdnaming.GRPCResolver{Client: cli}
	b := grpc.RoundRobin(r)

	conn, err := grpc.Dial("myService", grpc.WithBalancer(b), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	ticker := time.NewTicker(1 * time.Second)
	for t := range ticker.C {
		client := pb.NewGreeterClient(conn)
		resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "world " + strconv.Itoa(t.Second())})
		if err == nil {
			fmt.Printf("%v: Reply is %s\n", t, resp.Message)
		}
	}

}
