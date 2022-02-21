package discord_notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mroth/shopmon/internal/shopify"
)

type DiscordNotifier struct {
	webhook string
}

type payload struct {
	Username  string  `json:"username,omitempty"`   // Overrides the predefined username of the webhook
	AvatarURL string  `json:"avatar_url,omitempty"` // Overrides the predefined avatar of the webhook
	Content   string  `json:"content,omitempty"`    // Text message. Up to 2000 characters.
	Embeds    []embed `json:"embeds,omitempty"`     // Embed message.
}

type embed struct {
	Title       string    `json:"title,omitempty"`
	URL         string    `json:"url,omitempty"`
	Description string    `json:"description,omitempty"`
	Color       string    `json:"color,omitempty"`
	Author      author    `json:"author,omitempty"`
	Fields      []field   `json:"fields,omitempty"`
	Thumbnail   thumbnail `json:"thumbnail,omitempty"`
	Image       image     `json:"image,omitempty"`
	Footer      footer    `json:"footer,omitempty"`
}

type author struct {
	Name    string `json:"name,omitempty"`
	URL     string `json:"url,omitempty"`
	IconURL string `json:"icon_url,omitempty"`
}

type field struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}

type thumbnail struct {
	URL string `json:"url,omitempty"`
}

type image struct {
	URL string `json:"url,omitempty"`
}

type footer struct {
	Text    string `json:"text,omitempty"`
	IconURL string `json:"icon_url,omitempty"`
}

func postWebhookContext(ctx context.Context, webhook string, payload *payload) error {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhook, &body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %v", res.StatusCode)
	}

	return nil
}

func New(webhook string) DiscordNotifier {
	return DiscordNotifier{webhook: webhook}
}

func (n DiscordNotifier) Notify(storeDomain string, p shopify.ProductDetails) error {
	return n.NotifyWithContext(context.Background(), storeDomain, p)
}

func (n DiscordNotifier) NotifyWithContext(ctx context.Context, storeDomain string, p shopify.ProductDetails) error {
	msg := fmt.Sprintf("**%v** is now in stock! https://%s%s", p.Title, storeDomain, p.URL)
	return postWebhookContext(ctx, n.webhook, &payload{
		Content: msg,
		// Embeds: []embed{
		// 	{
		// 		Title:       p.Title,
		// 		URL:         fmt.Sprintf("https://%s%s", storeDomain, p.URL),
		// 		Description: p.Description,
		// 		Thumbnail:   thumbnail{URL: "https:" + p.FeaturedImage},
		// 	},
		// },
	})
}
