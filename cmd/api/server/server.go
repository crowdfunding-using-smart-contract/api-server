package server

import (
	"context"
	"errors"
	"fmt"
	"fund-o/api-server/cmd/worker"
	"fund-o/api-server/pkg/mail"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hibiken/asynq"

	"fund-o/api-server/internal/datasource"
	"fund-o/api-server/internal/datasource/repository"
	"fund-o/api-server/internal/http/handler"
	"fund-o/api-server/internal/http/middleware"
	"fund-o/api-server/internal/usecase"
	"fund-o/api-server/pkg/token"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

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
	log.Info().Msg("Starting listening for HTTP requests...")
	go func() {
		log.Info().Msgf("Server listening at http://%s:%d", server.config.APP_HOST, server.config.APP_PORT)
		log.Info().Msg("Starting listening for HTTP requests completed")
		if err := server.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("Failed to listen and serve")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Info().Msg("Shutting down server...")

	log.Info().Msg("Unregistering datasource...")
	if err := server.datasource.Close(); err != nil {
		return fmt.Errorf("error when close datasources: %v", err)
	}
	log.Info().Msg("Unregistering datasource completed")

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	if err := server.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("error when shutdown server: %v", err)
	}

	<-ctx.Done()
	log.Info().Msg("Timeout of 1 second")
	log.Info().Msg("Shutting down server completed")
	return nil
}

func inject(config *ApiServerConfig, datasource datasource.Datasource) *gin.Engine {
	// Makers
	jwtMaker, err := token.NewJWTMaker(config.JWT_SECRET_KEY)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create JWT maker")
	}

	// Repositories
	transactionRepository := repository.NewTransactionRepository(datasource.GetSqlDB())
	userRepository := repository.NewUserRepository(datasource.GetSqlDB())
	sessionRepository := repository.NewSessionRepository(datasource.GetSqlDB())
	projectRepository := repository.NewProjectRepository(datasource.GetSqlDB())
	projectCategoryRepository := repository.NewProjectCategoryRepository(datasource.GetSqlDB())
	verifyEmailRepository := repository.NewVerifyEmailRepository(datasource.GetSqlDB())
	forumRepository := repository.NewForumRepository(datasource.GetSqlDB())

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
	forumUseCase := usecase.NewForumUseCase(&usecase.ForumUseCaseOptions{
		ForumRepository: forumRepository,
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
		UserUseCase:        userUseCase,
		SessionUseCase:     sessionUseCase,
		VerifyEmailUseCase: verifyEmailUseCase,
		TokenMaker:         jwtMaker,
		TaskDistributor:    taskDistributor,
	})
	userHandler := handler.NewUserHandler(&handler.UserHandlerOptions{
		UserUseCase: userUseCase,
	})
	projectHandler := handler.NewProjectHandler(&handler.ProjectHandlerOptions{
		ProjectUseCase:         projectUseCase,
		UserUseCase:            userUseCase,
		ProjectCategoryUseCase: projectCategoryUseCase,
	})
	forumHandler := handler.NewForumHandler(&handler.ForumHandlerOptions{
		ForumUseCase: forumUseCase,
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
		authRoute.GET("/verify-email", authHandler.VerifyEmail)
		authRoute.POST("/send-verify-email", authHandler.SendVerifyEmail)
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
	postRoute := routeV1.Group("/posts")
	{
		postRoute.GET("", forumHandler.ListPosts)
		postRoute.POST("", authMiddleware, forumHandler.CreatePost)
		postRoute.GET("/:id", forumHandler.GetPostByID)
		postRoute.POST("/:id/comments", authMiddleware, forumHandler.CreateComment)
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
	mailer := mail.NewGmailSender(&gmailOptions)
	taskProcessor := worker.NewRedisTaskProcessor(&worker.RedisTaskProcessorOptions{
		RedisOptions: redisOptions,
		Mailer:       mailer,
		UseCases:     useCases,
	})

	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}
}
