package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Publisher struct {
	bot    *tgbotapi.BotAPI
	chatID int64
}

func (p Publisher) Publish(ctx context.Context, msg string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		_, err := p.bot.Send(tgbotapi.NewMessage(p.chatID, msg))
		if err != nil {
			return fmt.Errorf("send message: %w", err)
		}
		return nil
	}
}

func NewPublisher(token string, chatID int64) (*Publisher, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("telegram bot: %w", err)
	}
	return &Publisher{bot: bot, chatID: chatID}, nil
}
