package log

import (
	"os"

	"github.com/evalphobia/logrus_sentry"
	"github.com/getsentry/raven-go"
	"github.com/mhd7966/VulEQ/configs"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func LogInit() {
	Log = logrus.New()

	config := configs.Cfg.Log

	Log.SetFormatter(&logrus.TextFormatter{
		DisableColors:   false,
		TimestampFormat: "2006-Jan-02 15:04:05",
	})

	if config.OutputType == "stdout" {
		Log.SetOutput(os.Stdout)

	} else if config.OutputType == "file" {
		file, err := os.OpenFile(config.OutputAdd, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			Log.Fatal(err)
		}
		Log.SetOutput(file)
	} else {
		Log.SetOutput(os.Stdout)

	}

	logLevel, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	Log.SetLevel(logLevel)

	client, err := raven.New(configs.Cfg.Sentry.DSN)
	if err != nil {
		Log.Fatal(err)
	}

	hook, err := logrus_sentry.NewWithClientSentryHook(client, []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
	})

	if err == nil {
		Log.Hooks.Add(hook)
	}

	// err = raven.SetDSN(configs.Cfg.Sentry.DSN)
	// if err != nil {
	// 	logrus.Fatalf("%+v", errors.Wrap(err, "Can't set up raven"))
	// }

	// hook := sentryhook.New(nil)                        // will use raven.DefaultClient, or provide custom client
	// // hook.SetAsync()                   // async (non-bloking) hook for errors
	// hook.SetSync(logrus.ErrorLevel, logrus.PanicLevel, logrus.FatalLevel) // sync (blocking) for fatal stuff

	// logrus.AddHook(hook)

}
