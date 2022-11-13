package helpers

import (
	"context"
	"time"
)

/* GetTimeoutCtx Returns a timeout context with the specified seconds */
func GetTimeoutCtx(seconds string) (context.Context, context.CancelFunc) {
	timeout, _ := time.ParseDuration(seconds)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	return ctx, cancel
}
