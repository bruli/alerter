package message

import "context"

//go:generate go tool moq -out zmock_publisher.go . Publisher
type Publisher interface {
	Publish(ctx context.Context, msg string) error
}
