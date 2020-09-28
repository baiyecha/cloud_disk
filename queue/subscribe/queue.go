package subscribe

import (
	"context"
	"github.com/baiyecha/cloud_disk/pkg/pubsub"
	"github.com/baiyecha/cloud_disk/server"
)

func StartSubQueue(svr *server.Server) {
	ctx := context.Background()
	sub := pubsub.NewSub(svr.RedisClient, svr.Logger, svr.Conf.QueueNum)
	sub.RegisterSub()
	sub.Sub(ctx)
}
