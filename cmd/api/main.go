package main

import (
	"fund-o/api-server/cmd/api/server"
	"fund-o/api-server/config"
	"fund-o/api-server/internal/datasource"
	"fund-o/api-server/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func init() {
	logger.InitLogger()
}

func main() {
	appConfig, err := config.LoadAppConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load app config")
	}

	datasources := datasource.NewDatasourceContext(&appConfig.DatasourceConfig)

	gin.SetMode(appConfig.GinMode)
	s := server.NewApiServer(&appConfig.ApiServerConfig, datasources)
	if err := s.Start(); err != nil {
		log.Fatal().Err(err).Msg("Failed to start API server")
	}
}
