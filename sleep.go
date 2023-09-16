package ao3

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
)

// Code taken from
// https://github.com/chromedp/chromedp/blob/v0.9.2/chromedp.go#L742
func Sleep(d time.Duration) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		return sleepContext(ctx, d)
	})
}

func sleepContext(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	select {
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
