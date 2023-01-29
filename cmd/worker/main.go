package main

import (
	"context"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/hibiken/asynq"
	"github.com/mhd7966/VulEQ/configs"
	"github.com/mhd7966/VulEQ/connections"
	"github.com/mhd7966/VulEQ/jobs"
	"github.com/mhd7966/VulEQ/log"
	"github.com/mhd7966/VulEQ/services/sonarqube"
	"github.com/sirupsen/logrus"
)

func main() {

	configs.SetConfigs()
	log.LogInit()
	connections.ConnectDB()
	connections.ConnectRedis()
	sonarqube.Ping()

	if !configs.Cfg.Debug {
		err := sentry.Init(sentry.ClientOptions{})
		if err != nil {
			panic(err)
		}
		defer sentry.Flush(2 * time.Second)
	}

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: configs.Cfg.Redis.Addr},
		asynq.Config{
			Concurrency: 2,
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				sentry.ConfigureScope(func(scope *sentry.Scope) {
					scope.SetContext("Task", map[string]interface{}{
						"Payload": task.Payload(),
						"Type":    task.Type(),
					})
				})
				sentry.CaptureException(err)
			}),
		})

	mux := asynq.NewServeMux()
	log.Log.Info("Server Mux Create Succesfully :)")
	mux.HandleFunc(jobs.TypeCreateProject, jobs.HandleCreateProjectTask)
	mux.HandleFunc(jobs.TypeScanProject, jobs.HandleScanProjectTask)

	if err := srv.Run(mux); err != nil {
		log.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("There was a problem running the server")
	}
	log.Log.Info("Server Mux Run Succesfully :)")
}
