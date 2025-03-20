package notifier

import (
	"context"
	"fmt"
	"os"
)

type Notifier struct {
	webhookURL string
}

func NewNotifier() *Notifier {
	return &Notifier{
		webhookURL: os.Getenv("DISCORD_WEBHOOK_URL"),
	}
}

func (n *Notifier) Notify(ctx context.Context, message string) error {
	if n.webhookURL == "" {
		return fmt.Errorf("DISCORD_WEBHOOK_URL is not defined")
	}

	// TODO: Implement Discord notification
	fmt.Printf("Notification: %s\n", message)
	return nil
}
