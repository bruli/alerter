package memory_test

import (
	"context"
	"testing"
	"time"

	"github.com/bruli/alerter/internal/infra/memory"
	"github.com/stretchr/testify/require"
)

func TestCache_Set(t *testing.T) {
	t.Run("Given a Cache with running ttl job", func(t *testing.T) {
		t.Run(`when a value is set but it's expired,
		then Get method should return false'`, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			c := memory.NewCache(200 * time.Millisecond)
			go c.RunTTL(ctx, 500*time.Millisecond)
			value := "test"
			c.Set(value)
			time.Sleep(time.Second)
			require.False(t, c.Exists(value))
		})
		t.Run(`when a value is set but is alive,
		then Get method should return true`, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			c := memory.NewCache(500 * time.Millisecond)
			go c.RunTTL(ctx, 200*time.Millisecond)
			value := "test"
			c.Set(value)
			time.Sleep(100 * time.Millisecond)
			require.True(t, c.Exists(value))
		})
	})
}
