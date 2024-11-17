package usecase

import (
	"context"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/patyukin/mbs-pkg/pkg/model"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

func (u *UseCase) handleStartCommand(ctx context.Context, message *tgbotapi.Message) {
	args := message.CommandArguments()

	if args == "" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Welcome! Please provide a valid invite link to start.")
		if _, err := u.bot.API.Send(msg); err != nil {
			log.Error().Msgf("Error sending message: %v", err)
			return
		}

		log.Info().Msgf("Sent message to chat %d", message.Chat.ID)
		return
	}

	code, err := uuid.Parse(args)
	if err != nil {
		log.Error().Msgf("Error parsing invite code: %v", err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Invalid invite link.")
		if _, err = u.bot.API.Send(msg); err != nil {
			log.Error().Msgf("Error sending message: %v", err)
			return
		}

		log.Info().Msgf("Sent message to chat %d", message.Chat.ID)
		return
	}

	payload := model.AuthSignUpConfirmCode{
		Code:              code.String(),
		ChatID:            message.Chat.ID,
		UserTelegramLogin: message.From.UserName,
		UserTelegramID:    message.From.ID,
	}

	bytes, err := json.Marshal(payload)
	if err != nil {
		log.Error().Msgf("Error marshaling payload: %v", err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Error: "+err.Error())
		if _, err = u.bot.API.Send(msg); err != nil {
			log.Error().Msgf("Error sending message: %v", err)
			return
		}

		log.Info().Msgf("Sent message to chat %d", message.Chat.ID)
		return
	}

	err = u.rbt.PublishNotifySignUpConfirmCode(ctx, bytes, amqp.Table{})
	if err != nil {
		log.Error().Msgf("Error consuming message: %v", err)

		msg := tgbotapi.NewMessage(message.Chat.ID, "Error: "+err.Error())
		if _, err = u.bot.API.Send(msg); err != nil {
			log.Error().Msgf("Error sending message: %v", err)
			return
		}

		log.Info().Msgf("Sent message to chat %d", message.Chat.ID)
		return
	}
}
