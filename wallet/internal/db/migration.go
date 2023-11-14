package db

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"gorm.io/gorm"
	"log"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func MigrateSchemas(db *gorm.DB, dbName string) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("can't get sql.DB from gorm.DB: ", err)
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		log.Fatal("can't get driver: ", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://wallet/db/migrations",
		dbName, driver)
	if err != nil {
		log.Fatal("can't create migration instance: ", err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		version, dirty, _ := m.Version()
		if dirty {
			m.Force(int(version - 1))
		}
		log.Fatal("can't migrate schemas: ", err)
	}
}
