package dao

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"ppt/pg"
	"sync"
	"time"
)

var (
	pgDB *gorm.DB
	once sync.Once
)

func InitPg(cfg *pg.PgConfig) error {
	once.Do(func() {
		dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.SSLMode)
		pg := postgres.New(postgres.Config{DSN: dsn})
		db, err := gorm.Open(pg, &gorm.Config{})
		if err != nil {
			panic("Failed to connect to pg: " + err.Error())
		}
		sqlDB, err := db.DB()
		if err != nil {
			panic("Failed to connect to pg: " + err.Error())
		}
		sqlDB.SetMaxOpenConns(20)
		sqlDB.SetMaxIdleConns(5)
		sqlDB.SetConnMaxLifetime(time.Second * 10)
		if err = sqlDB.Ping(); err != nil {
			panic("Failed to connect to pg: " + err.Error())
		}
		pgDB = db
	})
	return nil
}

func ClosePg() {
	if pgDB != nil {
		sqlDB, err := pgDB.DB()
		if err != nil {
			fmt.Println("Failed to close pgDB: " + err.Error())
			return
		}
		_ = sqlDB.Close()
		pgDB = nil
	}
}

func GetPgDB() *gorm.DB {
	return pgDB
}
