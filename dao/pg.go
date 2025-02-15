package dao

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"ppt/pg"
	"sync"
	"time"
)

var (
	pgDB    *gorm.DB
	pgxPool *pgxpool.Pool
	once    sync.Once
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
		pgxPool, err = initPgxPool(dsn)
		if err != nil {
			panic("Failed to connect to pg: " + err.Error())
		}
	})
	return nil
}

func initPgxPool(dsn string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	// 连接池配置
	config.MinConns = 10
	config.MaxConns = 20
	config.MaxConnLifetime = time.Minute * 10
	config.MaxConnIdleTime = time.Minute * 5
	config.HealthCheckPeriod = time.Second * 30

	pool, err := pgxpool.NewWithConfig(pg.Ctx, config)
	if err != nil {
		return nil, err
	}
	if err = pool.Ping(pg.Ctx); err != nil {
		return nil, err
	}
	return pool, nil
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
