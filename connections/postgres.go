package connections

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/mhd7966/VulEQ/configs"
	"github.com/mhd7966/VulEQ/log"
	"github.com/sirupsen/logrus"
)

var DB *sql.DB

func ConnectDB() error {

	config := configs.Cfg.Database

	var err error

	DB, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Pass, config.Name))

	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Connect to DB Failed !!!!!")
		return err
	}

	err = DB.Ping()
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Ping DB Have Error !!!!!")
		return err
	}

	return nil

}
