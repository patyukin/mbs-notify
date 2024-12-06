package usecase

import (
	"context"

	"github.com/patyukin/mbs-notify/internal/telegram"
	amqp "github.com/rabbitmq/amqp091-go"
)

type KafkaProducer interface {
	PublishRegistrationSolution(ctx context.Context, value []byte) error
}

type UseCase struct {
	bot *telegram.Bot
	kfk KafkaProducer
}

func New(bot *telegram.Bot, kfk KafkaProducer) *UseCase {
	return &UseCase{
		bot: bot,
		kfk: kfk,
	}
}

func (u *UseCase) HandleAuthLogin(amqp.Delivery) error {
	return nil
}
