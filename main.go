package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"

	slack_notifier "github.com/mroth/shopmon/internal/notifiers/slack"
	"github.com/mroth/shopmon/internal/shopify"
)

type config struct {
	Domain         string   `env:"DOMAIN,notEmpty"`
	ProductHandles []string `env:"HANDLES,notEmpty" envSeparator:","`
	Rate           uint     `env:"RATE" envDefault:"60"`
	SlackWebhook   string   `env:"SLACK_WEBHOOK"`
	// DiscordWebhook string   `env:"DISCORD_WEBHOOK"`
}

type Notifier interface {
	Notify(productName, productURL string) error
}

type LogNotifier struct{}

func (n LogNotifier) Notify(productName, productURL string) error {
	log.Printf("🏪 %v is available! %v\n", productName, productURL)
	return nil
}

func main() {
	err := godotenv.Load()
	if err == nil {
		log.Println("INFO: loaded environment variables from local .env config file")
	}

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("%+v\n", err)
	}

	var notifiers []Notifier
	notifiers = append(notifiers, LogNotifier{})
	if cfg.SlackWebhook != "" {
		log.Println("INFO: configuring Slack notifier")
		notifiers = append(notifiers, slack_notifier.New(cfg.SlackWebhook))
	}

	ticker := time.NewTicker(time.Second * time.Duration(cfg.Rate))
	defer ticker.Stop()

	interuptC := make(chan os.Signal, 1)
	signal.Notify(interuptC, os.Interrupt)

	for {
		select {
		case <-interuptC:
			log.Println("INFO: received interrupt signal, shutting down...")
			os.Exit(0)
		case <-ticker.C:
			for _, handle := range cfg.ProductHandles {
				d, err := shopify.FetchProductDetails(cfg.Domain, handle)
				if err != nil {
					log.Printf("ERROR: %+v\n", err)
				} else {
					log.Printf("Checked %v: available %v\n", d.Title, d.Available)
					if d.Available {
						fullURL := fmt.Sprintf("https://%v/products/%v", cfg.Domain, d.Handle)
						for _, n := range notifiers {
							err := n.Notify(d.Title, fullURL)
							if err != nil {
								log.Printf("NOTIFICATION ERROR: %+v\n", err)
							}
						}
					}
				}
			}
		}
	}
}
