package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/patyukin/mbs-pkg/pkg/model"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

func (u *UseCase) ConsumeTelegramMessageQueue(ctx context.Context, msg amqp.Delivery) error {
	select {
	case <-ctx.Done():
		log.Error().Msgf("Context is done before processing: %v", ctx.Err())
		return ctx.Err()
	default:
	}

	var message model.SimpleTelegramMessage
	if err := json.Unmarshal(msg.Body, &message); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	log.Debug().Msgf("Received message: %+v", string(msg.Body))

	msgConfig := tgbotapi.NewMessage(message.ChatID, message.Message)
	if _, err := u.bot.API.Send(msgConfig); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	log.Info().Msgf("Sent message to chat %d", message.ChatID)

	return nil
}
