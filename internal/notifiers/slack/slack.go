package slack_notifier

import (
	"context"
	"fmt"

	"github.com/mroth/shopmon/internal/shopify"
	"github.com/slack-go/slack"
)

type SlackNotifier struct {
	webhook string
}

func New(webhook string) SlackNotifier {
	return SlackNotifier{webhook: webhook}
}

func (n SlackNotifier) Notify(storeDomain string, p shopify.ProductDetails) error {
	return n.NotifyWithContext(context.Background(), storeDomain, p)
}

func (n SlackNotifier) NotifyWithContext(ctx context.Context, storeDomain string, p shopify.ProductDetails) error {
	msg := fmt.Sprintf("*%v* is now in stock!\nhttps://%s%s", p.Title, storeDomain, p.URL)
	return slack.PostWebhookContext(ctx, n.webhook, &slack.WebhookMessage{
		Username:  "ShopMon Bot",
		IconEmoji: "convenience_store",
		// IconURL:         "",
		// Channel:         "",
		Text: msg,
		Blocks: &slack.Blocks{
			BlockSet: []slack.Block{
				slack.NewSectionBlock(
					slack.NewTextBlockObject("mrkdwn", msg, false, false),
					[]*slack.TextBlockObject{
						// slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Price:* %d", p.Price), false, false),
					},
					slack.NewAccessory(slack.NewImageBlockElement("https:"+p.FeaturedImage, p.Handle)),
				),
			},
		},
	})
}
