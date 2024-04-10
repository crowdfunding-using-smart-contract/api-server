package server

import (
	"context"
	"errors"
	"fmt"
	"fund-o/api-server/cmd/worker"
	"fund-o/api-server/cmd/ws"
	"fund-o/api-server/config"
	"fund-o/api-server/pkg/mail"
	"fund-o/api-server/pkg/uploader"
	"github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
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
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"

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
	config     *config.ApiServerConfig
	datasource datasource.Datasource
}

func NewApiServer(config *config.ApiServerConfig, datasource datasource.Datasource) ApiServer {
	router := inject(config, datasource)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Host, config.Port),
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
		log.Info().Msgf("Server listening at http://%s:%d", server.config.Host, server.config.Port)
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

func inject(config *config.ApiServerConfig, datasource datasource.Datasource) *gin.Engine {
	// Makers
	jwtMaker, err := token.NewJWTMaker(config.JwtSecretKey)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create JWT maker")
	}

	imageUploader, err := uploader.NewS3Store(&uploader.S3StoreConfig{
		Region:             config.AwsRegion,
		Bucket:             config.AwsBucketName,
		AwsAccessKeyID:     config.AwsAccessKeyID,
		AwsSecretAccessKey: config.AwsSecretAccessKey,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create image uploader")
	}

	// Repositories
	transactionRepository := repository.NewTransactionRepository(datasource.GetSqlDB())
	userRepository := repository.NewUserRepository(datasource.GetSqlDB())
	sessionRepository := repository.NewSessionRepository(datasource.GetSqlDB())
	projectRepository := repository.NewProjectRepository(datasource.GetSqlDB())
	projectCategoryRepository := repository.NewProjectCategoryRepository(datasource.GetSqlDB())
	verifyEmailRepository := repository.NewVerifyEmailRepository(datasource.GetSqlDB())
	forumRepository := repository.NewForumRepository(datasource.GetSqlDB())
	channelRepository := repository.NewChannelRepository(datasource.GetSqlDB())
	messageRepository := repository.NewMessageRepository(datasource.GetSqlDB())

	// UseCases
	transactionUseCase := usecase.NewTransactionUseCase(&usecase.TransactionUseCaseOptions{
		TransactionRepository: transactionRepository,
	})
	userUseCase := usecase.NewUserUseCase(&usecase.UserUseCaseOptions{
		UserRepository: userRepository,
		ImageUploader:  imageUploader,
	})
	sessionUseCase := usecase.NewSessionUseCase(&usecase.SessionUseCaseOptions{
		SessionRepository: sessionRepository,
	})
	projectUseCase := usecase.NewProjectUseCase(&usecase.ProjectUseCaseOptions{
		ProjectRepository: projectRepository,
		ImageUploader:     imageUploader,
	})
	projectCategoryUseCase := usecase.NewProjectCategoryUseCase(&usecase.ProjectCategoryUseCaseOptions{
		ProjectCategoryRepository: projectCategoryRepository,
	})
	verifyEmailUseCase := usecase.NewVerifyEmailUseCase(&usecase.VerifyEmailUseCaseOptions{
		VerifyEmailRepository: verifyEmailRepository,
	})
	forumUseCase := usecase.NewForumUseCase(&usecase.ForumUseCaseOptions{
		ForumRepository: forumRepository,
		ImageUploader:   imageUploader,
	})
	channelUsecase := usecase.NewChannelUsecase(&usecase.ChannelUsecaseOptions{
		ChannelRepository: channelRepository,
	})
	messageUseCase := usecase.NewMessageUsecase(&usecase.MessageUsecaseOptions{
		MessageRepository: messageRepository,
		ImageUploader:     imageUploader,
	})

	// Task Processor
	redisOptions := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOptions)
	gmailOptions := mail.GmailSenderOptions{
		Name:              config.EmailSenderName,
		FromEmailAddress:  config.EmailSenderAddress,
		FromEmailPassword: config.EmailSenderPassword,
	}
	go runTaskProcessor(redisOptions, gmailOptions, &worker.TaskProcessorUseCaseOptions{
		UserUseCase:        userUseCase,
		VerifyEmailUseCase: verifyEmailUseCase,
	})

	// Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.RedisAddress,
	})

	// Websocket
	hub := ws.NewWebsocketHub(&ws.Config{
		Redis: redisClient,
	})
	go hub.Run()
	socketService := ws.NewSocketService(&ws.SocketServiceConfig{
		Hub: hub,
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
	chatHandler := handler.NewChatHandler(&handler.ChatHandlerOptions{
		ChannelUsecase: channelUsecase,
		MessageUsecase: messageUseCase,
		SocketService:  socketService,
	})

	authMiddleware := middleware.AuthMiddleware(jwtMaker)

	router := gin.New()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{config.CorsAllowedOrigin},
		AllowCredentials: config.CorsAllowedCredentials,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
	})
	router.Use(c)

	router.Use(middleware.RequestLogger())
	router.Use(middleware.ResponseLogger())
	router.Use(registerRateLimiter(redisClient))

	router.GET("/ws", authMiddleware, func(c *gin.Context) {
		ws.ServeWs(hub, c)
	})

	routeV1 := router.Group(config.PathPrefix)

	docs.SwaggerInfo.BasePath = config.PathPrefix
	initSwaggerDocs(routeV1)

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
		userRoute.PATCH("/:id", authMiddleware, userHandler.UpdateUser)
	}
	projectRoute := routeV1.Group("/projects")
	{
		projectRoute.GET("", projectHandler.ListProjects)
		projectRoute.POST("", authMiddleware, projectHandler.CreateProject)
		projectRoute.GET("/:id", projectHandler.GetProjectByID)
		projectRoute.GET("/me", authMiddleware, projectHandler.GetOwnProjects)
		projectRoute.GET("/recommendation", projectHandler.GetRecommendProjects)
		projectRoute.GET("/categories", projectHandler.ListProjectCategories)
		projectRoute.POST("/:id/ratings", authMiddleware, projectHandler.CreateProjectRating)
		projectRoute.GET("/:id/ratings/verify", authMiddleware, projectHandler.VerifyProjectRating)
		projectRoute.POST("/:id/contribute", authMiddleware, projectHandler.ContributeProject)
		projectRoute.GET("/backed", authMiddleware, projectHandler.GetBackedProject)
	}
	postRoute := routeV1.Group("/posts")
	{
		postRoute.GET("", forumHandler.ListPosts)
		postRoute.POST("", authMiddleware, forumHandler.CreatePost)
		postRoute.GET("/:id", forumHandler.GetPostByID)
		postRoute.POST("/:id/comments", authMiddleware, forumHandler.CreateComment)
		postRoute.POST("/upload", authMiddleware, forumHandler.UploadImage)
	}
	commentRoute := routeV1.Group("/comments")
	{
		commentRoute.POST("/:id/replies", authMiddleware, forumHandler.CreateReply)
	}
	channelRoute := routeV1.Group("/channels")
	{
		channelRoute.GET("/me", authMiddleware, chatHandler.GetOwnChannels)
		channelRoute.GET("/:id", authMiddleware, chatHandler.GetOrCreateChannel)
		channelRoute.POST("/:id/messages", authMiddleware, chatHandler.SendMessage)
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

func registerRateLimiter(redisClient *redis.Client) gin.HandlerFunc {
	rate, err := limiter.NewRateFromFormatted("5-S")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create rate limiter")
	}

	store, err := sredis.NewStoreWithOptions(redisClient, limiter.StoreOptions{
		Prefix: "limiter_fundo",
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create rate limiter store")
	}

	rateLimitMiddleware := mgin.NewMiddleware(limiter.New(store, rate))
	return rateLimitMiddleware
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
