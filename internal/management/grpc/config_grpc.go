package grpc

import (
	"sync"
	"sync/atomic"
	"time"
	"zmd-gacha/internal/models"

	pb "zmd-gacha/proto"

	"github.com/google/uuid"
)

type ConfigUpdateEvent struct {
	Version uint64
	Config  models.GachaPoolConfig
}

type ConfigHub struct {
	mu          sync.RWMutex
	subscribers map[string]chan ConfigUpdateEvent
	latest      atomic.Pointer[ConfigUpdateEvent] // 新订阅可先收到最新快照
}

func NewConfigHub() *ConfigHub {
	return &ConfigHub{subscribers: make(map[string]chan ConfigUpdateEvent)}
}

func (h *ConfigHub) Subscribe(id string) (<-chan ConfigUpdateEvent, func()) {
	ch := make(chan ConfigUpdateEvent, 64)

	h.mu.Lock()
	h.subscribers[id] = ch
	h.mu.Unlock()

	// 新连接先推一次最新快照（如果有）
	if p := h.latest.Load(); p != nil {
		select {
		case ch <- *p:
		default:
		}
	}

	cancel := func() {
		h.mu.Lock()
		if c, ok := h.subscribers[id]; ok {
			delete(h.subscribers, id)
			close(c)
		}
		h.mu.Unlock()
	}
	return ch, cancel
}

func (h *ConfigHub) Publish(evt ConfigUpdateEvent) {
	h.latest.Store(&evt)

	h.mu.RLock()
	defer h.mu.RUnlock()

	for id, ch := range h.subscribers {
		select {
		case ch <- evt:
		default:
			// 慢消费者队列满：可记录日志并选择丢弃旧消息或断开该客户端
			_ = id
		}
	}
}

type ConfigSyncServer struct {
	pb.UnimplementedConfigSyncServiceServer
	hub *ConfigHub
}

func NewConfigSyncServer(hub *ConfigHub) *ConfigSyncServer {
	return &ConfigSyncServer{hub: hub}
}

func (s *ConfigSyncServer) SubscribePoolConfig(
	req *pb.SubscribeRequest,
	stream pb.ConfigSyncService_SubscribePoolConfigServer,
) error {
	subID := req.GetServerId()
	if subID == "" {
		subID = uuid.NewString()
	}

	ch, cancel := s.hub.Subscribe(subID)
	defer cancel()

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case evt, ok := <-ch:
			if !ok {
				return nil
			}

			rsp := &pb.ConfigUpdate{
				Version:         evt.Version,
				UpdatedAtUnixMs: time.Now().UnixMilli(),
				Config: &pb.PoolConfig{
					PoolId:               uint32(evt.Config.PoolID),
					SRankBaseRate:        evt.Config.SRankBaseRate,
					ARankBaseRate:        evt.Config.ARankBaseRate,
					AGuaranteeInterval:   int32(evt.Config.AGuaranteeInterval),
					SPityStart:           int32(evt.Config.SPityStart),
					SPityStep:            evt.Config.SPityStep,
					SPityEnd:             int32(evt.Config.SPityEnd),
					LimitPity:            int32(evt.Config.LimitPity),
					LimitRateWhenS:       evt.Config.LimitRateWhenS,
					MaxLimitedCharacters: int32(evt.Config.MaxLimitedCharacters),
				},
			}
			if err := stream.Send(rsp); err != nil {
				return err
			}
		}
	}
}
