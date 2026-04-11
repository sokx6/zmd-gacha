package service

import (
	"context"
	"errors"
	"io"
	"log"
	"time"
	"zmd-gacha/internal/models"
	pb "zmd-gacha/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func StartConfigWatcher(ctx context.Context, managementAddr string, serverID string, gs *GachaService) {
	if managementAddr == "" {
		managementAddr = "127.0.0.1:9090"
	}

	backoff := time.Second
	for {
		err := consumeConfigStream(ctx, managementAddr, serverID, gs)
		if err != nil && !errors.Is(err, context.Canceled) {
			log.Printf("config watcher disconnected: %v", err)
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(backoff):
			if backoff < 30*time.Second {
				backoff *= 2
			}
		}
	}
}

func consumeConfigStream(ctx context.Context, managementAddr string, serverID string, gs *GachaService) error {
	conn, err := grpc.DialContext(ctx, managementAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewConfigSyncServiceClient(conn)
	stream, err := client.SubscribePoolConfig(ctx, &pb.SubscribeRequest{ServerId: serverID})
	if err != nil {
		return err
	}

	for {
		msg, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		if msg.GetConfig() == nil {
			continue
		}

		cfg := models.GachaPoolConfig{
			PoolID:               uint(msg.GetConfig().GetPoolId()),
			SRankBaseRate:        msg.GetConfig().GetSRankBaseRate(),
			ARankBaseRate:        msg.GetConfig().GetARankBaseRate(),
			AGuaranteeInterval:   int(msg.GetConfig().GetAGuaranteeInterval()),
			SPityStart:           int(msg.GetConfig().GetSPityStart()),
			SPityStep:            msg.GetConfig().GetSPityStep(),
			SPityEnd:             int(msg.GetConfig().GetSPityEnd()),
			LimitPity:            int(msg.GetConfig().GetLimitPity()),
			LimitRateWhenS:       msg.GetConfig().GetLimitRateWhenS(),
			MaxLimitedCharacters: int(msg.GetConfig().GetMaxLimitedCharacters()),
		}
		gs.ApplyConfigUpdate(msg.GetVersion(), cfg)
	}
}
