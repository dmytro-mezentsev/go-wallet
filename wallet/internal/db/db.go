package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"net/url"
	"wallet.com/wallet/wallet/internal/config"
)

func DBConnection(conf config.DbConf) *gorm.DB {
	dsn := url.URL{
		User:     url.UserPassword(conf.User, conf.Password),
		Scheme:   "postgres",
		Host:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Path:     conf.DBName,
		RawQuery: (&url.Values{"sslmode": []string{"disable"}}).Encode(),
	}
	db, err := gorm.Open(postgres.Open(dsn.String()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("can't connect to db: ", err)
	}
	return db

}
