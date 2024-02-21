package driver

import (
	"errors"
	"fmt"

	"fund-o/api-server/internal/datasource/driver/seeds"
	"fund-o/api-server/internal/entity"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type SqlDBConfig struct {
	SQL_HOST     string `mapstructure:"SQL_HOST"`
	SQL_USERNAME string `mapstructure:"SQL_USERNAME"`
	SQL_PASSWORD string `mapstructure:"SQL_PASSWORD"`
	SQL_PORT     int    `mapstructure:"SQL_PORT"`
	SQL_DATABASE string `mapstructure:"SQL_DATABASE"`
}

type SQLContext interface {
	Connect() error
	Disconnect() error
	DB() *gorm.DB
}

type sqlContext struct {
	dsn string
	db  *gorm.DB
}

var logger *log.Entry

func NewSQLContext(config *SqlDBConfig) SQLContext {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Bangkok",
		config.SQL_HOST,
		config.SQL_USERNAME,
		config.SQL_PASSWORD,
		config.SQL_DATABASE,
		config.SQL_PORT,
	)

	logger = log.WithFields(log.Fields{
		"dsn": fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s",
			config.SQL_USERNAME,
			"*****",
			config.SQL_HOST,
			config.SQL_PORT,
			config.SQL_DATABASE,
		),
	})

	return &sqlContext{dsn: dsn}
}

func (sql *sqlContext) Connect() error {
	logger.Info("Connecting to SQL database...")

	if sql.dsn == "" {
		return fmt.Errorf("failed to connect to SQL database: DSN is empty")
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: sql.dsn,
	}), &gorm.Config{
		// Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err != nil {
		return err
	}
	logger.Info("Connecting to SQL database completed")

	sql.db = db

	if err := sql.autoMigrateUp(); err != nil {
		return err
	}

	return nil
}

func (sql *sqlContext) Disconnect() error {
	logger.Info("Disconnecting from SQL database...")

	if sql.db != nil {
		db, err := sql.db.DB()
		if err != nil {
			return err
		}
		err = db.Close()
		if err != nil {
			return err
		}
	}

	logger.Info("Disconnecting from SQL database completed")

	return nil
}

func (sql *sqlContext) DB() *gorm.DB {
	return sql.db
}

func (sql *sqlContext) autoMigrateUp() error {
	db := sql.db
	if err := db.AutoMigrate(
		&entity.Transaction{},
		&entity.User{},
		&entity.Session{},
		&entity.Project{},
		&entity.ProjectCategory{},
		&entity.ProjectSubCategory{},
	); err != nil {
		return err
	}

	if err := db.AutoMigrate(&entity.ProjectCategory{}); err == nil && db.Migrator().HasTable(&entity.ProjectCategory{}) {
		if err := db.First(&entity.ProjectCategory{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			for _, category := range seeds.ProjectCategorySeed {
				if result := db.Create(&category); result.Error != nil {
					return result.Error
				}
			}
		}
	}

	return nil
}
