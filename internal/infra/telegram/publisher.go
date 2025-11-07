package telegram

import (
	"context"
	"fmt"

	"github.com/bruli/alerter/internal/domain/message"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Publisher struct {
	bot    *tgbotapi.BotAPI
	chatID int64
}

func (p Publisher) Publish(ctx context.Context, m *message.Message) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		_, err := p.bot.Send(tgbotapi.NewMessage(p.chatID, fmt.Sprintf("⚠️ failed ping. Resource: %q", m.Resource())))
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
