package middleware

// func init() {
// 	raven.SetDSN("

// https://examplePublicKey@o0.ingest.sentry.io/0")
// }

import (
	"time"

	sentryfiber "github.com/aldy505/sentry-fiber"
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"github.com/mhd7966/VulEQ/configs"
)

func InitApiSentry(app fiber.Router) {

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              configs.Cfg.Sentry.DSN,
		AttachStacktrace: true,
	})
	if err != nil {
		panic(err)
	}
	app.Use(sentryfiber.New(sentryfiber.Options{}))
	defer sentry.Flush(2 * time.Second)

}
