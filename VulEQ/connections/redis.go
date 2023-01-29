package connections

import (
	"github.com/abr-ooo/VulEQ/configs"
	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
	// log "github.com/sirupsen/logrus"
)

var RedisClient *asynq.Client

func ConnectRedis() {

	RedisClient = asynq.NewClient(asynq.RedisClientOpt{Addr: configs.Cfg.Redis.Addr})


}
