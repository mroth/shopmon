package slack_notifier

import (
	"context"
	"fmt"

	"github.com/slack-go/slack"
)

type SlackNotifier struct {
	webhook string
}

func New(webhook string) SlackNotifier {
	return SlackNotifier{webhook: webhook}
}

func (n SlackNotifier) Notify(productName, productURL string) error {
	return n.NotifyWithContext(context.Background(), productName, productURL)
}

func (n SlackNotifier) NotifyWithContext(ctx context.Context, productName, productURL string) error {
	return slack.PostWebhookContext(ctx, n.webhook, &slack.WebhookMessage{
		Username:  "ShopMon Bot",
		IconEmoji: "convenience_store",
		// IconURL:         "",
		// Channel:         "",
		// ThreadTimestamp: "",
		Text: fmt.Sprintf("*%v* is now in stock!\n%v", productName, productURL),
		// Attachments:     []slack.Attachment{},
		// Parse: "",
		// Blocks: &slack.Blocks{
		// 	BlockSet: []slack.Block{},
		// },
		// ResponseType:    "",
		// ReplaceOriginal: false,
		// DeleteOriginal:  false,
	})
}
