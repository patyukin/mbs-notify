package usecase

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

const startCommand = "start"

func (u *UseCase) StartTelegramBot(ctx context.Context) {
	up := tgbotapi.NewUpdate(0)
	up.Timeout = 60

	updates := u.bot.API.GetUpdatesChan(up)

	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		log.Info().Msgf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		switch update.Message.Command() {
		case startCommand:
			u.handleStartCommand(ctx, update.Message)
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command. Use /help to see the list of available commands.")
			if _, err := u.bot.API.Send(msg); err != nil {
				log.Error().Msgf("Error sending message: %v", err)
				return
			}

			log.Info().Msgf("Sent message to chat %d", update.Message.Chat.ID)
		}
	}
}
