package datasource

import (
	"fund-o/api-server/internal/datasource/driver"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type DatasourceConfig struct {
	driver.SqlDBConfig
}

type Datasource interface {
	GetSqlDB() *gorm.DB
	Close() error
}

type datasource struct {
	sql driver.SQLContext
}

func NewDatasourceContext(config *DatasourceConfig) Datasource {
	sqlDBContext := driver.NewSQLContext(&config.SqlDBConfig)

	err := sqlDBContext.Connect()
	if err != nil {
		log.Panic().Err(err).Msg("Failed to connect to SQL database")
	}

	return &datasource{
		sql: sqlDBContext,
	}
}

func (ds *datasource) GetSqlDB() *gorm.DB {
	return ds.sql.DB()
}

func (ds *datasource) Close() error {
	err := ds.sql.Disconnect()
	return err
}
