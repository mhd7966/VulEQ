package configs

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
)

var Cfg Config

type Config struct {
	Debug    bool `env:"DEBUG" env-default:"False"`
	Database struct {
		Port string `env:"DB_PORT" env-default:"5432"`
		Host string `env:"DB_HOST" env-default:"db"`
		Name string `env:"DB_NAME" env-default:"postgres"`
		User string `env:"DB_USER" env-default:"admin"`
		Pass string `env:"DB_PASS" env-default:"postgres_password"`
	}
	///sonarscanner/sonar-scanner-4.6.2.2472-linux/bin/sonar-scanner
	SonarQube struct {
		Host        string `env:"SONAR_HOST" env-default:"sonarqube:9000"`
		User        string `env:"SONAR_USER" env-default:"admin"`
		Pass        string `env:"SONAR_PASS" env-default:"sonar_password"`
		Auth        string `env:"Authorization" env-default:"authorization_sonar"`
		ScannerPath string `env:"SONAR_SCANNER_PATH" env-default:"/opt/sonar-scanner/bin/sonar-scanner"`
		TokenPrefix string `env:"SONAR_TOKEN_PREFIX" env-default:"token_"`
		KeyPrefix   string `env:"SONAR_KEY_PREFIX" env-default:"key_"`
	}
	Git struct {
		ClonePath string `env:"STORE_CLONE_PATH" env-default:"/git_clone/"`
	}
	Redis struct {
		Addr string `env:"REDIS_ADD" env-default:"redis:6379"`
	}
	Log struct {
		LogLevel   string `env:"LOG_LEVEL" env-default:"debug"`
		OutputType string `env:"LOG_OUTPUT_TYPE" env-default:"stdout"`
		OutputAdd  string `env:"LOG_FILE_Add" env-default:"/log.txt"`
	}
	Auth struct {
		Host string `env:"AUTH_HOST" env-default:"authorization_host_address(del)"`
	}
	Sentry struct {
		DSN string `env:"SENTRY_DSN" env-default:"sentry_dsn_address"`
		Level string `env:"SENTRY_LEVEL" env-default:"error"`
	}
}

func SetConfigs() {

	if _, err := os.Stat(".env"); err == nil {
		cleanenv.ReadConfig(".env", &Cfg)
		logrus.Info("Set config from .env file")
	} else {
		cleanenv.ReadEnv(&Cfg)
		logrus.Info("Set config from Config struct values")
	}

}
