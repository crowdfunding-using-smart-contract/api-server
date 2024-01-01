package main

import (
	"github.com/danyouknowme/gin-gorm-boilerplate/cmd/api/server"
	"github.com/danyouknowme/gin-gorm-boilerplate/config"
	"github.com/danyouknowme/gin-gorm-boilerplate/internal/datasource"
	"github.com/danyouknowme/gin-gorm-boilerplate/pkg/logger"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func init() {
	logger.InitLogger()
}

func main() {
	appConfig, err := config.LoadAppConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	datasources := datasource.NewDatasourceContext(&appConfig.DatasourceConfig)

	gin.SetMode(appConfig.GIN_MODE)
	server := server.NewApiServer(&appConfig.ApiServerConfig, datasources)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
