package usecase

import (
	"context"
	"github.com/patyukin/mbs-notify/internal/telegram"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ interface {
	PublishDQLMessage(ctx context.Context, body []byte) error
	PublishNotifySignUpConfirmCode(ctx context.Context, body []byte, headers amqp.Table) error
	PublishAuthSignUpResultMessage(ctx context.Context, body []byte, headers amqp.Table) error
	PublishAuthSignInCode(ctx context.Context, body []byte, headers amqp.Table) error
}

type UseCase struct {
	bot *telegram.Bot
	rbt RabbitMQ
}

func New(bot *telegram.Bot, rbt RabbitMQ) *UseCase {
	return &UseCase{
		bot: bot,
		rbt: rbt,
	}
}

func (u *UseCase) HandleAuthLogin(amqp.Delivery) error {
	return nil
}
