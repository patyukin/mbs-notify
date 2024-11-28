package usecase

import (
	"context"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/patyukin/mbs-pkg/pkg/model"
	"github.com/rs/zerolog/log"
)

func (u *UseCase) sendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := u.bot.API.Send(msg); err != nil {
		log.Error().Msgf("Error sending message to chat %d: %v", chatID, err)
		return err
	}

	log.Info().Msgf("Sent message to chat %d", chatID)
	return nil
}

func (u *UseCase) handleStartCommand(ctx context.Context, message *tgbotapi.Message) {
	args := message.CommandArguments()

	if args == "" {
		log.Error().Msg("No invite link provided")
		if err := u.sendMessage(message.Chat.ID, "Welcome! Please provide a valid invite link to start."); err != nil {
			return
		}
		return
	}

	code, err := uuid.Parse(args)
	if err != nil {
		log.Error().Msgf("Error parsing invite code: %v", err)
		if sendErr := u.sendMessage(message.Chat.ID, "Invalid invite link."); sendErr != nil {
			return
		}
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
		if sendErr := u.sendMessage(message.Chat.ID, "Error was encountered. Please try again."); sendErr != nil {
			return
		}
		return
	}

	err = u.kfk.PublishRegistrationSolution(ctx, bytes)
	if err != nil {
		log.Error().Msgf("Error consuming message: %v", err)
		if sendErr := u.sendMessage(message.Chat.ID, "Error was encountered. Please try again."); sendErr != nil {
			return
		}
		return
	}
}
