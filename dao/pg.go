package dao

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"ppt/config"
	"ppt/log"
	"ppt/nacos/wrapper"
	"sync"
	"time"
)

var (
	PgDB    *gorm.DB
	pgxPool *pgxpool.Pool
	pgOnce  sync.Once
)

func InitPg(cfg *wrapper.PgConfig) error {
	var err error
	pgOnce.Do(func() {
		dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", cfg.UserName, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)
		PgDB, err = initPgGorm(dsn)
		if err != nil {
			panic(err)
		}

		pgxPool, err = initPgxPool(dsn)
		if err != nil {
			panic("Failed to connect to pg: " + err.Error())
		}
	})
	return nil
}

func initPgGorm(dsn string) (*gorm.DB, error) {
	pg := postgres.New(postgres.Config{DSN: dsn})
	db, err := gorm.Open(pg, &gorm.Config{})
	if err != nil {
		log.Error("initPgGorm error", zap.String("dsn", dsn), zap.Error(err))
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Error("initPgGorm failed to connect pg", zap.String("dsn", dsn), zap.Error(err))
		return nil, err
	}
	// 微服务适当降低配置
	//sqlDB.SetMaxOpenConns(20)
	//sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(3)
	// 定时清理连接资源
	sqlDB.SetConnMaxLifetime(time.Hour * 1)
	// 及时释放多余连接
	sqlDB.SetConnMaxIdleTime(time.Minute * 3)
	if err = sqlDB.Ping(); err != nil {
		log.Error("initPgGorm failed to ping pg", zap.String("dsn", dsn), zap.Error(err))
		return nil, err
	}
	db = db.Debug()
	if config.Env == "prod" {
		db.Logger = logger.Default.LogMode(logger.Silent)
	}
	return db, nil
}

func initPgxPool(dsn string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Error("initPgxPool parse config error", zap.String("dsn", dsn), zap.Error(err))
		return nil, err
	}
	// 连接池配置
	config.MinConns = 10
	config.MaxConns = 20
	config.MaxConnLifetime = time.Minute * 10
	config.MaxConnIdleTime = time.Minute * 5
	config.HealthCheckPeriod = time.Second * 30

	pool, err := pgxpool.NewWithConfig(Ctx, config)
	if err != nil {
		log.Error("initPgxPool new pool with config error", zap.String("dsn", dsn), zap.Any("pool_config", *config), zap.Error(err))
		return nil, err
	}
	if err = pool.Ping(Ctx); err != nil {
		log.Error("initPgxPool pool ping error", zap.String("dsn", dsn), zap.Any("pool_config", *config), zap.Error(err))
		return nil, err
	}
	return pool, nil
}

func ClosePg() {
	if PgDB != nil {
		sqlDB, err := PgDB.DB()
		if err != nil {
			fmt.Println("Failed to close pgDB: " + err.Error())
			return
		}
		_ = sqlDB.Close()
		PgDB = nil
	}
	if pgxPool != nil {
		pgxPool.Close()
		pgxPool = nil
	}
}
