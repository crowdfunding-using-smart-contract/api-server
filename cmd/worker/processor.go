package worker

import (
	"context"
	"fund-o/api-server/internal/usecase"
	"fund-o/api-server/pkg/mail"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server   *asynq.Server
	mailer   mail.EmailSender
	useCases *TaskProcessorUseCaseOptions
	logger   *Logger
}

type RedisTaskProcessorOptions struct {
	RedisOptions asynq.RedisClientOpt
	Mailer       mail.EmailSender
	UseCases     *TaskProcessorUseCaseOptions
}

type TaskProcessorUseCaseOptions struct {
	UserUseCase        usecase.UserUseCase
	VerifyEmailUseCase usecase.VerifyEmailUseCase
}

func NewRedisTaskProcessor(options *RedisTaskProcessorOptions) TaskProcessor {
	logger := NewWorkerLogger("processor")
	redis.SetLogger(logger)

	server := asynq.NewServer(
		asynq.RedisClientOpt{},
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				logger.log.Error().Err(err).
					Str("type", task.Type()).
					Bytes("payload", task.Payload()).
					Msg("process task failed")
			}),
			Logger: logger,
		},
	)

	return &RedisTaskProcessor{
		server: server,
		mailer: options.Mailer,
		useCases: &TaskProcessorUseCaseOptions{
			UserUseCase:        options.UseCases.UserUseCase,
			VerifyEmailUseCase: options.UseCases.VerifyEmailUseCase,
		},
		logger: logger,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	log := processor.logger.log

	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)

	log.Info().Msg("Starting task processor...")
	go func() {
		err := processor.server.Start(mux)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to start task processor")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Info().Msg("Shutting down task processor...")

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	processor.server.Shutdown()
	<-ctx.Done()

	return nil
}
