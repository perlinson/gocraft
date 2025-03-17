package main

import (
	"context"
	"crypto/tls"
	"log"
	"strconv"
	"time"

	playerPB "github.com/perlinson/gocraft/internal/proto/player"

	blockPB "github.com/perlinson/gocraft/internal/proto/block"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type GRPCClient struct {
	conn         *grpc.ClientConn
	blockClient  blockPB.BlockServiceClient
	playerClient playerPB.PlayerServiceClient
	// 连接池和重试机制待补充
}

// 实现方块更新方法
func (c *GRPCClient) UpdateBlock(id Vec3, w int) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req := &blockPB.UpdateBlockRequest{
		X: int32(id.X),
		Y: int32(id.Y),
		Z: int32(id.Z),
		W: int32(w),
	}

	resp, err := c.blockClient.UpdateBlock(ctx, req)
	if err != nil {
		log.Printf("方块更新失败: %v", err)
		return
	}

	if resp.Version != "" {
		store.UpdateChunkVersion(id.Chunkid(), resp.Version)
	}
}

// 实现玩家状态同步方法
func (c *GRPCClient) UpdatePlayerState(state PlayerState) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req := &playerPB.UpdateStateRequest{
		State: &playerPB.PlayerState{
			X:  float32(state.X),
			Y:  float32(state.Y),
			Z:  float32(state.Z),
			Rx: float32(state.Rx),
			Ry: float32(state.Ry),
		},
	}

	resp, err := c.playerClient.UpdateState(ctx, req)
	if err != nil {
		log.Printf("玩家状态更新失败: %v", err)
		return
	}

	for id, player := range resp.Players {
		// 将string类型的id转换为int32后传入UpdateOrAdd
		idInt, _ := strconv.ParseInt(id, 10, 32)
		game.playerRender.UpdateOrAdd(int32(idInt), player)
	}
}

// func InitGRPCClient() (*GRPCClient, error) {
// 	if *serverAddr == "" {
// 		return nil, nil
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	conn, err := grpc.DialContext(ctx, *serverAddr,
// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// 		grpc.WithBlock(),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &GRPCClient{
// 		conn:         conn,
// 		blockClient:  blockPB.NewBlockServiceClient(conn),
// 		playerClient: playerPB.NewPlayerServiceClient(conn),
// 	}, nil
// }

func InitGRPCClient() (*GRPCClient, error) {
	if *serverAddr == "" {
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 创建 TLS 配置
	creds := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: true, // 开发环境可用
	})

	conn, err := grpc.DialContext(ctx, *serverAddr,
		grpc.WithTransportCredentials(creds), // 使用 TLS
		grpc.WithBlock(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		return nil, err
	}

	return &GRPCClient{
		conn:         conn,
		blockClient:  blockPB.NewBlockServiceClient(conn),
		playerClient: playerPB.NewPlayerServiceClient(conn),
	}, nil
}

// StreamChunkUpdates 实现区块流式更新处理，包含错误重试和连接保持机制
func (c *GRPCClient) StreamChunkUpdates(ctx context.Context, id Vec3, callback func(bid Vec3, w int)) {
	retryCount := 0
	maxRetries := 3

	for {
		stream, err := c.blockClient.StreamChunk(ctx, &blockPB.ChunkRequest{
			P:       int32(id.X),
			Q:       int32(id.Z),
			Version: "",
		})
		if err != nil {
			log.Printf("建立区块流失败: %v", err)
			if retryCount >= maxRetries {
				return
			}
			time.Sleep(time.Duration(retryCount+1) * time.Second)
			retryCount++
			continue
		}

		retryCount = 0
		go c.handleChunkStream(stream, id, callback)
		return
	}
}

func (c *GRPCClient) handleChunkStream(stream blockPB.BlockService_StreamChunkClient, id Vec3, callback func(bid Vec3, w int)) {
	for {
		update, err := stream.Recv()
		if err != nil {
			log.Printf("流接收错误: %v", err)
			return
		}

		// 处理区块中的所有方块更新
		for i := 0; i < len(update.Blocks); i += 4 {
			bid := Vec3{
				X: int(update.Blocks[i]),
				Y: int(update.Blocks[i+1]),
				Z: int(update.Blocks[i+2]),
			}
			w := int(update.Blocks[i+3])

			// 检查方块坐标的有效性
			if bid.Chunkid() != id {
				log.Printf("跳过无效方块更新: 方块ID %v 不在区块 %v 中", bid, id)
				continue
			}

			store.UpdateBlock(bid, w)
			game.blockRender.DirtyChunk(bid.Chunkid())

			if callback != nil {
				callback(bid, w)
			}
		}

		if update.Version != "" {
			store.UpdateChunkVersion(id, update.Version)
		}
	}
}
