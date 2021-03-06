package main

// 降级的使用

import (
	"context"
	"fmt"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-micro/v2/web"
	"gomicro_note/p19/grpc_client/models"
	"gomicro_note/p19/grpc_client/routers"
	"gomicro_note/p19/grpc_client/wrappers"
)

type logWrapper struct {
	client.Client
}

func (l *logWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	md, _ := metadata.FromContext(ctx)
	fmt.Printf("[Log Wrapper] ctx: %v service: %s method: %s\n", md, req.Service(), req.Endpoint())
	return l.Client.Call(ctx, req, rsp)
}

func newLogWrapper(c client.Client) client.Client {
	return &logWrapper{c}
}

func main() {

	etcdReg := etcd.NewRegistry(
		registry.Addrs("127.0.0.1:2379"),
	)

	myService := micro.NewService(
		micro.Name("ProdService.client"),
		micro.WrapClient(newLogWrapper),
		micro.WrapClient(wrappers.NewProdsWrapper),
	)
	prodService := models.NewProdService("ProdService", myService.Client())

	service := web.NewService(
		web.Name("ProdService.client"),
		web.Address(":9000"),
		web.Handler(routers.InitRouter(prodService)),
		web.Registry(etcdReg),
	)

	service.Init()
	service.Run()
}
