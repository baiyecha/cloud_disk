package wrapper

import (
	"context"
	"github.com/baiyecha/cloud_disk/pkg/pubsub"
	"github.com/baiyecha/cloud_disk/store/db_store"
	"github.com/jinzhu/gorm"
)

type DB struct {
	sub pubsub.SubQueue
	db  *gorm.DB
}

func (g *DB) Channel() string {
	return g.sub.Channel()
}

func (g *DB) Process(ctx context.Context, message string) {
	g.sub.Process(db_store.NewDBContext(ctx, g.db), message)
}

func NewDB(sub pubsub.SubQueue, db *gorm.DB) pubsub.SubQueue {
	return &DB{sub: sub, db: db}
}
