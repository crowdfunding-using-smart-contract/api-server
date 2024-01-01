package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danyouknowme/gin-gorm-boilerplate/internal/datasource"
	"github.com/danyouknowme/gin-gorm-boilerplate/internal/datasource/repository"
	"github.com/danyouknowme/gin-gorm-boilerplate/internal/http/handler"
	"github.com/danyouknowme/gin-gorm-boilerplate/internal/http/middleware"
	"github.com/danyouknowme/gin-gorm-boilerplate/internal/usecase"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"

	cors "github.com/rs/cors/wrapper/gin"

	docs "github.com/danyouknowme/gin-gorm-boilerplate/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type ApiServer interface {
	Start() error
	HttpServer() *http.Server
}

type apiServer struct {
	httpServer  *http.Server
	config      *ApiServerConfig
	datasources datasource.Datasource
}

func NewApiServer(config *ApiServerConfig, datasources datasource.Datasource) ApiServer {
	router := inject(config, datasources)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.APP_HOST, config.APP_PORT),
		Handler: router,
	}

	return &apiServer{
		httpServer:  server,
		config:      config,
		datasources: datasources,
	}
}

func (server *apiServer) HttpServer() *http.Server {
	return server.httpServer
}

func (server *apiServer) Start() error {
	logger.Info("Starting listening for HTTP requests...")
	go func() {
		logger.Infof("Server listening at http://%s:%d", server.config.APP_HOST, server.config.APP_PORT)
		logger.Info("Starting listening for HTTP requests completed")
		if err := server.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Info("Shutting down server...")

	logger.Info("Unregistering datasources...")
	if err := server.datasources.Close(); err != nil {
		return fmt.Errorf("error when close datasources: %v", err)
	}
	logger.Info("Unregistering datasources completed")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := server.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("error when shutdown server: %v", err)
	}

	<-ctx.Done()
	logger.Info("Timeout of 1 second")
	logger.Info("Shutting down server completed")
	return nil
}

func inject(config *ApiServerConfig, datasources datasource.Datasource) *gin.Engine {
	// Repositories
	transactionRepository := repository.NewTransactionRepository(datasources.GetSqlDB())

	// Usecases
	transactionUsecase := usecase.NewTransactionUsecase(&usecase.TransactionUsecaseOptions{
		TransactionRepository: transactionRepository,
	})

	// Handlers
	transactionHandler := handler.NewTransactionHandler(&handler.TransactionHandlerOptions{
		TransactionUsecase: transactionUsecase,
	})

	router := gin.New()
	routeV1 := router.Group(config.APP_PATH_PREFIX)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{config.APP_CORS_ALLOWED_ORIGIN},
		AllowCredentials: config.APP_CORS_ALLOWED_CREDENTIALS,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
	})
	routeV1.Use(c)

	routeV1.Use(middleware.RequestLogger())
	routeV1.Use(middleware.ResponseLogger())

	docs.SwaggerInfo.BasePath = config.APP_PATH_PREFIX
	initSwaggerDocs(routeV1)

	// Routes
	routeV1.GET("/hello", handler.GetHelloMessage)
	transactionRoute := routeV1.Group("/transactions")
	{
		transactionRoute.GET("", transactionHandler.ListTransactions)
		transactionRoute.GET("/:id", transactionHandler.GetTransaction)
		transactionRoute.POST("", transactionHandler.CreateTransaction)
	}

	return router
}

// @title FundO API
// @version 1.0
// @description This is a sample server server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @schemes http https

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func initSwaggerDocs(router *gin.RouterGroup) {
	router.GET("/openapi/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
