package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"fund-o/api-server/internal/entity"
	"fund-o/api-server/pkg/helper"
	log "github.com/sirupsen/logrus"

	"github.com/hibiken/asynq"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
	Email string `json:"email"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Errorf("failed to marshal task payload: %+v", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		log.Errorf("failed to enqueue task: %+v", err)
	}

	log.WithFields(log.Fields{
		"type":      task.Type(),
		"payload":   string(task.Payload()),
		"queue":     info.Queue,
		"max_retry": info.MaxRetry,
	}).Info("enqueued task")
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(_ context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	user, err := processor.useCases.UserUseCase.GetUserByEmail(payload.Email)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	verifyEmail, err := processor.useCases.VerifyEmailUseCase.CreateVerifyEmail(&entity.VerifyEmailCreatePayload{
		Email:      user.Email,
		SecretCode: helper.RandomString(32),
	})
	if err != nil {
		return fmt.Errorf("failed to create verify email: %w", err)
	}

	subject := "Welcome to Simple Bank"
	// TODO: replace this URL with an environment variable that points to a front-end page
	verifyUrl := fmt.Sprintf("http://localhost:5173/v1/verify_email?email_id=%s&secret_code=%s", verifyEmail.ID, verifyEmail.SecretCode)
	content := fmt.Sprintf(`Hello %s,<br/>
	Thank you for registering with us!<br/>
	Please <a href="%s">click here</a> to verify your email address.<br/>
	`, user.FullName, verifyUrl)
	to := []string{user.Email}

	err = processor.mailer.SendEmail(subject, content, to, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}

	log.WithFields(log.Fields{
		"type":    task.Type(),
		"payload": string(task.Payload()),
		"email":   user.Email,
	}).Info("processed task")
	return nil
}
