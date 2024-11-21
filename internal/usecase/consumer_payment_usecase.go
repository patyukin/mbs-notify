package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/patyukin/mbs-pkg/pkg/model"
	"github.com/patyukin/mbs-pkg/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

func (u *UseCase) PaymentConsumeHandler(ctx context.Context, msg amqp.Delivery) error {
	select {
	case <-ctx.Done():
		log.Error().Msgf("Context is done before processing in PaymentConsumeHandler: %v", ctx.Err())
		return ctx.Err()
	default:
	}

	switch msg.RoutingKey {
	case rabbitmq.AccountCreationRouteKey, rabbitmq.PaymentExecutionInitiateRouteKey:
		var message model.SimpleTelegramMessage
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			return fmt.Errorf("failed to unmarshal message for routing key '%s': %w", msg.RoutingKey, err)
		}

		msgConfig := tgbotapi.NewMessage(message.ChatID, message.Message)
		if _, err := u.bot.API.Send(msgConfig); err != nil {
			return fmt.Errorf("failed to send message for routing key '%s': %w", msg.RoutingKey, err)
		}

		log.Info().Msgf("Sent message to chat %d, from %s", message.ChatID, msg.RoutingKey)

	default:
		return fmt.Errorf("unknown routing key: %s", msg.RoutingKey)
	}

	log.Info().Msg("consumer done")

	return nil
}
