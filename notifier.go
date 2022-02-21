package main

import (
	"context"
	"log"

	"github.com/mroth/shopmon/internal/shopify"
)

type Notifier interface {
	NotifyWithContext(ctx context.Context, storeDomain string, p shopify.ProductDetails) error
	Notify(storeDomain string, p shopify.ProductDetails) error
}

// LogNotifier is a notifier that logs notifications to stdout.
type LogNotifier struct{}

func (n LogNotifier) NotifyWithContext(_ context.Context, storeDomain string, p shopify.ProductDetails) error {
	return n.Notify(storeDomain, p)
}

func (n LogNotifier) Notify(storeDomain string, p shopify.ProductDetails) error {
	log.Printf("🏪 %v is available! https://%s%s\n", p.Title, storeDomain, p.URL)
	return nil
}

var _ Notifier = LogNotifier{}
