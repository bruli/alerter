package message_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bruli/alerter/internal/domain/message"
	"github.com/stretchr/testify/require"
)

func TestPublish_Handle(t *testing.T) {
	readyMessage, err := message.NewMessage("test", "ok")
	require.NoError(t, err)
	failedMessage, err := message.NewMessage("test", message.FailedStatus)
	require.NoError(t, err)
	ctx := context.Background()
	type args struct {
		ctx context.Context
		m   *message.Message
	}
	tests := []struct {
		name        string
		args        args
		expectedErr error
		exits       bool
		expectedSetCalls,
		expectedRemoveCalls int
	}{
		{
			name: "with a failed message and resource already exist in cache, then it no set the value in cache and not send message",
			args: args{
				ctx: ctx,
				m:   failedMessage,
			},
			exits:            true,
			expectedSetCalls: 0,
		},
		{
			name: "with a failed message and resource not exist in cache but publisher returns and error, then it returns same error",
			args: args{
				ctx: ctx,
				m:   failedMessage,
			},
			exits:            false,
			expectedErr:      errors.New("test"),
			expectedSetCalls: 0,
		},
		{
			name: "with a failed message and resource not exist in cache, then it set in cache and returns send message",
			args: args{
				ctx: ctx,
				m:   failedMessage,
			},
			exits:            false,
			expectedSetCalls: 1,
		},
		{
			name: "with a ready message and resource exit exist in cache, then it remove cache and publish message",
			args: args{
				ctx: ctx,
				m:   readyMessage,
			},
			exits:               true,
			expectedRemoveCalls: 1,
		},
		{
			name: "with a ready message and resource exit exist in cache but publish returns an error, then it returns same error",
			args: args{
				ctx: ctx,
				m:   readyMessage,
			},
			exits:       true,
			expectedErr: errors.New("test"),
		},
	}
	for _, tt := range tests {
		t.Run(`Given a Publish service,
		when Handle method is called `+tt.name, func(t *testing.T) {
			pub := &message.PublisherMock{}
			pub.PublishFunc = func(_ context.Context, msg string) error {
				return tt.expectedErr
			}
			cache := &message.CacheMock{}
			cache.ExistsFunc = func(a string) bool {
				return tt.exits
			}
			cache.SetFunc = func(v string) {}
			cache.RemoveFunc = func(v string) {}

			p := message.NewPublish(pub, cache)
			err = p.Handle(tt.args.ctx, tt.args.m)
			if err != nil {
				require.Equal(t, tt.expectedErr, err)
				return
			}
			require.Len(t, cache.SetCalls(), tt.expectedSetCalls)
		})
	}
}
