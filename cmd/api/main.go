package main

import (
	"fund-o/api-server/cmd/api/server"
	"fund-o/api-server/config"
	"fund-o/api-server/internal/datasource"
	"fund-o/api-server/pkg/logger"

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
