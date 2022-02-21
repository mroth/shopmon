package main

import (
	"context"
	"log"
)

type Notifier interface {
	NotifyWithContext(ctx context.Context, productName, productURL string) error
	Notify(productName, productURL string) error
}

// LogNotifier is a notifier that logs notifications to stdout.
type LogNotifier struct{}

func (n LogNotifier) NotifyWithContext(_ context.Context, productName, productURL string) error {
	return n.Notify(productName, productURL)
}

func (n LogNotifier) Notify(productName, productURL string) error {
	log.Printf("🏪 %v is available! %v\n", productName, productURL)
	return nil
}
