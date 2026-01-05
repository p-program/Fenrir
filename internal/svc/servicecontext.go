package svc

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"zeusro.com/gotemplate/internal/config"
	"zeusro.com/gotemplate/model"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := initDB(c)
	return &ServiceContext{
		Config: c,
		DB:     db,
	}
}

func initDB(c config.Config) *gorm.DB {
	var db *gorm.DB
	var err error

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	switch c.Database.Type {
	case "mysql":
		db, err = gorm.Open(mysql.Open(c.Database.DSN), gormConfig)
	case "postgres":
		db, err = gorm.Open(postgres.Open(c.Database.DSN), gormConfig)
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(c.Database.DSN), gormConfig)
	default:
		// 默认使用 SQLite
		db, err = gorm.Open(sqlite.Open("restaurant.db"), gormConfig)
	}

	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// 自动迁移
	if err := autoMigrate(db); err != nil {
		panic("failed to migrate database: " + err.Error())
	}

	return db
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Wallet{},
		&model.Transaction{},
		&model.Plate{},
		&model.Food{},
		&model.Order{},
		&model.OrderItem{},
		&model.PlateDepot{},
		&model.Worker{},
		&model.ExceptionLog{},
		&model.GCProcessLog{},
	)
}
