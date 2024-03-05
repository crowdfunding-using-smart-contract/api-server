package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeTaskSendVerifyEmail(
		ctx context.Context,
		payload *PayloadSendVerifyEmail,
		opts ...asynq.Option,
	)
}

type RedisTaskDistributor struct {
	client *asynq.Client
	logger *Logger
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpt)
	logger := NewWorkerLogger("distributor")

	return &RedisTaskDistributor{
		client: client,
		logger: logger,
	}
}
