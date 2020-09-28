package wrapper

import (
	"context"
	"github.com/baiyecha/cloud_disk/pkg/pubsub"
	"github.com/baiyecha/cloud_disk/service"
)

type Service struct {
	sub     pubsub.SubQueue
	service service.Service
}

func (g *Service) Channel() string {
	return g.sub.Channel()
}

func (g *Service) Process(ctx context.Context, message string) {
	g.sub.Process(service.NewContext(ctx, g.service), message)
}

func NewService(sub pubsub.SubQueue, service service.Service) pubsub.SubQueue {
	return &Service{sub: sub, service: service}
}
