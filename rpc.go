package main

import (
	"context"
	"flag"
	"time"
)

var (
	serverAddr = flag.String("s", "ywnzjusupjyh.sealosgzg.site:443", "server address")
)

func InitClient() error {
	var err error
	grpcClient, err = InitGRPCClient()
	return err
}

var grpcClient *GRPCClient

func ClientFetchChunk(id Vec3, f func(bid Vec3, w int)) {
	if grpcClient == nil {
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	grpcClient.StreamChunkUpdates(ctx, id, f)
}

func ClientUpdateBlock(id Vec3, w int) {
	if grpcClient == nil {
		return
	}

	grpcClient.UpdateBlock(id, w)
}

func ClientUpdatePlayerState(state PlayerState) {
	if grpcClient == nil {
		return
	}
	grpcClient.UpdatePlayerState(state)
}

// CloseClient 关闭gRPC客户端连接
func CloseClient() {
	if grpcClient != nil && grpcClient.conn != nil {
		grpcClient.conn.Close()
		grpcClient = nil
	}
}
