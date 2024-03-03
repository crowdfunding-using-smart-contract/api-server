package worker

import (
	"context"
	"fund-o/api-server/internal/usecase"
	"fund-o/api-server/pkg/mail"
	log "github.com/sirupsen/logrus"

	"github.com/hibiken/asynq"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	Start() error
	Shutdown()
}

type RedisTaskProcessor struct {
	server             *asynq.Server
	userUseCase        usecase.UserUseCase
	verifyEmailUseCase usecase.VerifyEmailUseCase
	mailer             mail.EmailSender
}

type RedisTaskProcessorOptions struct {
	RedisOptions asynq.RedisClientOpt
	Mailer       mail.EmailSender
	usecase.UserUseCase
	usecase.VerifyEmailUseCase
}

func NewRedisTaskProcessor(options *RedisTaskProcessorOptions) TaskProcessor {
	logger := log.WithFields(log.Fields{
		"module": "task_processor",
	})

	server := asynq.NewServer(
		options.RedisOptions,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				logger.WithError(err).WithFields(log.Fields{
					"type":    task.Type(),
					"payload": string(task.Payload()),
				}).Error("process task failed")
			}),
			Logger: logger,
		},
	)

	return &RedisTaskProcessor{
		server:             server,
		mailer:             options.Mailer,
		userUseCase:        options.UserUseCase,
		verifyEmailUseCase: options.VerifyEmailUseCase,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	return processor.server.Start(mux)
}

func (processor *RedisTaskProcessor) Shutdown() {
	processor.server.Shutdown()
}
