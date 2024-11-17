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

func (u *UseCase) AuthConsumeHandler(ctx context.Context, msg amqp.Delivery) error {
	select {
	case <-ctx.Done():
		log.Error().Msgf("Context is done before processing: %v", ctx.Err())
		return ctx.Err()
	default:
	}

	switch msg.RoutingKey {
	case rabbitmq.AuthSignUpResultMessageRouteKey:
		var message model.AuthSignUpResultMessage
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			return fmt.Errorf("failed to unmarshal message: %w", err)
		}

		msgConfig := tgbotapi.NewMessage(message.ChatID, message.Message)
		if _, err := u.bot.API.Send(msgConfig); err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}

		log.Info().Msgf("Sent message to chat %d, from rabbitmq.AuthSignUpResultMessageRouteKey", message.ChatID)

	case rabbitmq.AuthSignInConfirmCodeRouteKey:
		var authSignInCode model.AuthSignInCode
		err := json.Unmarshal(msg.Body, &authSignInCode)
		if err != nil {
			return fmt.Errorf("failed to unmarshal message body: %w", err)
		}

		msgConfig := tgbotapi.NewMessage(authSignInCode.ChatID, authSignInCode.Code)
		if _, err = u.bot.API.Send(msgConfig); err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}

		log.Info().Msgf("Sent message to chat %d, from rabbitmq.AuthSignInConfirmCodeRouteKey", authSignInCode.ChatID)

	default:
		return fmt.Errorf("unknown routing key: %s", msg.RoutingKey)
	}

	log.Info().Msg("consumer done")

	return nil
}
