package connections

import (
	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
	"github.com/mhd7966/VulEQ/configs"
	// log "github.com/sirupsen/logrus"
)

var RedisClient *asynq.Client

func ConnectRedis() {

	RedisClient = asynq.NewClient(asynq.RedisClientOpt{Addr: configs.Cfg.Redis.Addr})

}
