package server

import (
	"context"
	"errors"
	"fmt"
	"fund-o/api-server/cmd/worker"
	"fund-o/api-server/pkg/mail"
	"github.com/hibiken/asynq"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fund-o/api-server/internal/datasource"
	"fund-o/api-server/internal/datasource/repository"
	"fund-o/api-server/internal/http/handler"
	"fund-o/api-server/internal/http/middleware"
	"fund-o/api-server/internal/usecase"
	"fund-o/api-server/pkg/token"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	cors "github.com/rs/cors/wrapper/gin"

	docs "fund-o/api-server/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type ApiServer interface {
	Start() error
	HttpServer() *http.Server
}

type apiServer struct {
	httpServer *http.Server
	config     *ApiServerConfig
	datasource datasource.Datasource
}

func NewApiServer(config *ApiServerConfig, datasource datasource.Datasource) ApiServer {
	router := inject(config, datasource)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.APP_HOST, config.APP_PORT),
		Handler: router,
	}

	return &apiServer{
		httpServer: server,
		config:     config,
		datasource: datasource,
	}
}

func (server *apiServer) HttpServer() *http.Server {
	return server.httpServer
}

func (server *apiServer) Start() error {
	log.Info("Starting listening for HTTP requests...")
	go func() {
		log.Infof("Server listening at http://%s:%d", server.config.APP_HOST, server.config.APP_PORT)
		log.Info("Starting listening for HTTP requests completed")
		if err := server.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Info("Shutting down server...")

	log.Info("Unregistering datasource...")
	if err := server.datasource.Close(); err != nil {
		return fmt.Errorf("error when close datasources: %v", err)
	}
	log.Info("Unregistering datasource completed")

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	if err := server.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("error when shutdown server: %v", err)
	}

	<-ctx.Done()
	log.Info("Timeout of 1 second")
	log.Info("Shutting down server completed")
	return nil
}

func inject(config *ApiServerConfig, datasource datasource.Datasource) *gin.Engine {
	// Makers
	jwtMaker, err := token.NewJWTMaker(config.JWT_SECRET_KEY)
	if err != nil {
		log.Fatalf("Failed to create JWT maker: %v", err)
	}

	// Repositories
	transactionRepository := repository.NewTransactionRepository(datasource.GetSqlDB())
	userRepository := repository.NewUserRepository(datasource.GetSqlDB())
	sessionRepository := repository.NewSessionRepository(datasource.GetSqlDB())
	projectRepository := repository.NewProjectRepository(datasource.GetSqlDB())
	projectCategoryRepository := repository.NewProjectCategoryRepository(datasource.GetSqlDB())
	verifyEmailRepository := repository.NewVerifyEmailRepository(datasource.GetSqlDB())

	// UseCases
	transactionUseCase := usecase.NewTransactionUseCase(&usecase.TransactionUseCaseOptions{
		TransactionRepository: transactionRepository,
	})
	userUseCase := usecase.NewUserUseCase(&usecase.UserUseCaseOptions{
		UserRepository: userRepository,
	})
	sessionUseCase := usecase.NewSessionUseCase(&usecase.SessionUseCaseOptions{
		SessionRepository: sessionRepository,
	})
	projectUseCase := usecase.NewProjectUseCase(&usecase.ProjectUseCaseOptions{
		ProjectRepository: projectRepository,
	})
	projectCategoryUseCase := usecase.NewProjectCategoryUseCase(&usecase.ProjectCategoryUseCaseOptions{
		ProjectCategoryRepository: projectCategoryRepository,
	})
	verifyEmailUseCase := usecase.NewVerifyEmailUseCase(&usecase.VerifyEmailUseCaseOptions{
		VerifyEmailRepository: verifyEmailRepository,
	})

	// Task Processor
	redisOptions := asynq.RedisClientOpt{
		Addr: config.REDIS_ADDRESS,
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOptions)
	gmailOptions := mail.GmailSenderOptions{
		Name:              config.EMAIL_SENDER_NAME,
		FromEmailAddress:  config.EMAIL_SENDER_ADDRESS,
		FromEmailPassword: config.EMAIL_SENDER_PASSWORD,
	}
	go runTaskProcessor(redisOptions, gmailOptions, &worker.TaskProcessorUseCaseOptions{
		UserUseCase:        userUseCase,
		VerifyEmailUseCase: verifyEmailUseCase,
	})

	// Handlers
	transactionHandler := handler.NewTransactionHandler(&handler.TransactionHandlerOptions{
		TransactionUseCase: transactionUseCase,
	})
	authHandler := handler.NewAuthHandler(&handler.AuthHandlerOptions{
		UserUseCase:     userUseCase,
		SessionUseCase:  sessionUseCase,
		TokenMaker:      jwtMaker,
		TaskDistributor: taskDistributor,
	})
	userHandler := handler.NewUserHandler(&handler.UserHandlerOptions{
		UserUseCase: userUseCase,
	})
	projectHandler := handler.NewProjectHandler(&handler.ProjectHandlerOptions{
		ProjectUseCase:         projectUseCase,
		UserUseCase:            userUseCase,
		ProjectCategoryUseCase: projectCategoryUseCase,
	})

	router := gin.New()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{config.APP_CORS_ALLOWED_ORIGIN},
		AllowCredentials: config.APP_CORS_ALLOWED_CREDENTIALS,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
	})
	router.Use(c)

	router.Use(middleware.RequestLogger())
	router.Use(middleware.ResponseLogger())

	routeV1 := router.Group(config.APP_PATH_PREFIX)

	docs.SwaggerInfo.BasePath = config.APP_PATH_PREFIX
	initSwaggerDocs(routeV1)

	authMiddleware := middleware.AuthMiddleware(jwtMaker)

	// Routes
	routeV1.GET("/hello", handler.GetHelloMessage)
	transactionRoute := routeV1.Group("/transactions")
	{
		transactionRoute.GET("", transactionHandler.ListTransactions)
		transactionRoute.GET("/:id", transactionHandler.GetTransaction)
		transactionRoute.POST("", transactionHandler.CreateTransaction)
	}

	authRoute := routeV1.Group("/auth")
	{
		authRoute.POST("/register", authHandler.Register)
		authRoute.POST("/login", authHandler.Login)
		authRoute.POST("/renew-token", authHandler.RenewAccessToken)
		// authRoute.POST("/login-with-google", func(c *gin.Context) {
		// 	c.Redirect(http.StatusTemporaryRedirect, )
		// })
	}

	userRoute := routeV1.Group("/users")
	{
		userRoute.GET("/me", authMiddleware, userHandler.GetMe)
	}

	projectRoute := routeV1.Group("/projects")
	{
		projectRoute.POST("", authMiddleware, projectHandler.CreateProject)
		projectRoute.GET("/own", authMiddleware, projectHandler.GetOwnProjects)
		projectRoute.GET("/categories", projectHandler.ListProjectCategories)
	}

	return router
}

// @title FundO API
// @version 1.0
// @description This is a sample server.
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

func runTaskProcessor(
	redisOptions asynq.RedisClientOpt,
	gmailOptions mail.GmailSenderOptions,
	useCases *worker.TaskProcessorUseCaseOptions,
) {
	logger := log.WithFields(log.Fields{
		"module": "task_processor",
	})

	mailer := mail.NewGmailSender(&gmailOptions)
	taskProcessor := worker.NewRedisTaskProcessor(&worker.RedisTaskProcessorOptions{
		RedisOptions: redisOptions,
		Mailer:       mailer,
		UseCases:     useCases,
	})

	logger.Info("Starting task processor...")
	go func() {
		err := taskProcessor.Start()
		if err != nil {
			logger.Fatal("failed to start task processor")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Info("Shutting down task processor...")

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	taskProcessor.Shutdown()
	<-ctx.Done()
	log.Info("Shutting down task processor completed")
}
