package message_test

import (
	"testing"

	"github.com/bruli/alerter/internal/domain/message"
	"github.com/stretchr/testify/require"
)

func TestNewMessage(t *testing.T) {
	type args struct {
		resource string
		status   string
	}
	tests := []struct {
		name        string
		args        args
		expectedErr error
	}{
		{
			name:        "with an empty resource, then it returns an invalid resource error",
			args:        args{},
			expectedErr: message.ErrInvalidResource,
		},
		{
			name: "with an empty status, then it returns an invalid status error",
			args: args{
				resource: "bla",
				status:   "",
			},
			expectedErr: message.ErrInvalidStatus,
		},
		{
			name: "with all data, then it returns a valid message",
			args: args{
				resource: "bla",
				status:   message.FailedStatus,
			},
			expectedErr: message.ErrInvalidStatus,
		},
	}
	for _, tt := range tests {
		t.Run(`Given a NewMessage constructor,
		when is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := message.NewMessage(tt.args.resource, tt.args.status)
			if err != nil {
				require.ErrorIs(t, err, tt.expectedErr)
				return
			}
			require.Equal(t, tt.args.resource, got.Resource())
			require.Equal(t, tt.args.status, got.Status())
			require.True(t, got.IsFailed())
		})
	}
}
