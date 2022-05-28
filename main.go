package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"

	"github.com/mroth/shopmon/internal/notifiers/discord"
	"github.com/mroth/shopmon/internal/notifiers/slack"
	"github.com/mroth/shopmon/internal/shopify"
)

type config struct {
	Domain         string        `env:"DOMAIN,notEmpty"`                   // Shopify domain for store. (required)
	ProductHandles []string      `env:"HANDLES,notEmpty" envSeparator:","` // Product handles to check, comma separated. (required)
	Rate           uint          `env:"RATE" envDefault:"60"`              // How often to poll for products, in seconds. (default: 60)
	SlackWebhook   string        `env:"SLACK_WEBHOOK"`                     // Slack webhook URL to post notifications. (optional)
	DiscordWebhook string        `env:"DISCORD_WEBHOOK"`                   // Discord webhook URL to post notifications. (optional)
	FetchTimeout   time.Duration `env:"FETCH_TIMEOUT" envDefault:"10s"`    // Timeout for fetching product details. (default: 10s)
	NotifyTimeout  time.Duration `env:"NOTIFY_TIMEOUT" envDefault:"5s"`    // Timeout for posting a notification. (default: 5s)
}

func main() {
	// local development convenience, load .env file if exists
	err := godotenv.Load()
	if err == nil {
		log.Println("INFO: loaded environment variables from local .env config file")
	}

	// parse configuration from environment
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("%+v\n", err)
	}

	// setup a NotifyGroup according to configuration
	notifyGroup := setupNotifiers(cfg)
	for _, n := range notifyGroup.Notifiers {
		log.Printf("INFO: configured notifications via %v", reflect.TypeOf(n))
	}

	// ticker for update checks
	ticker := time.NewTicker(time.Second * time.Duration(cfg.Rate))
	defer ticker.Stop()

	// capture interrupt signals for graceful shutdown
	interuptC := make(chan os.Signal, 1)
	signal.Notify(interuptC, os.Interrupt)

	// primary event loop
	var wg sync.WaitGroup
	rootCtx, rootCancelF := context.WithCancel(context.Background())
	defer rootCancelF()

	for {
		select {
		case <-interuptC:
			log.Println("INFO: received interrupt signal, shutting down...")
			rootCancelF() // cancel in-flight checks and notifications
			wg.Wait()     // wait for completion
			os.Exit(0)
		case <-ticker.C:
			for _, handle := range cfg.ProductHandles {
				handle := handle

				wg.Add(1)
				go func() {
					defer wg.Done()
					ctx, cf := context.WithTimeout(rootCtx, cfg.FetchTimeout)
					defer cf()

					d, err := shopify.FetchProductDetails(ctx, cfg.Domain, handle)
					if err != nil {
						log.Printf("ERROR: %+v\n", err)
					} else {
						log.Printf("Checked %v: available %v\n", d.Title, d.Available)
						if d.Available {
							notifyGroup.Send(rootCtx, cfg.Domain, d)
						}
					}
				}()
			}
		}
	}
}

// NotifyGroup handles a collection of notifiers
type NotifyGroup struct {
	Notifiers []Notifier    // collection of configured notifiers
	Timeout   time.Duration // additional timeout restriction on notify
}

// Send concurrent notifications about ProductDetails from Shopify domain to all
// Notifiers in this NotifyGroup.  This method blocks until all notification
// goroutines have either completed successfully, error, or timeout.
//
// If a notification errors for any reason, it is logged to the global logger.
func (ng NotifyGroup) Send(ctx context.Context, domain string, d *shopify.ProductDetails) {
	var wg sync.WaitGroup
	for _, n := range ng.Notifiers {
		n := n
		wg.Add(1)
		go func() {
			defer wg.Done()

			ctx, cf := context.WithTimeout(ctx, ng.Timeout)
			defer cf()

			err := n.NotifyWithContext(ctx, domain, *d)
			if err != nil {
				log.Printf("NOTIFICATION ERROR: %+v\n", err)
			}
		}()
	}
	wg.Wait()
}

func setupNotifiers(cfg config) NotifyGroup {
	var notifiers []Notifier
	notifiers = append(notifiers, LogNotifier{}) // always include LogNotifier for now

	if cfg.SlackWebhook != "" {
		notifiers = append(notifiers, slack.New(cfg.SlackWebhook))
	}

	if cfg.DiscordWebhook != "" {
		notifiers = append(notifiers, discord.New(cfg.DiscordWebhook))
	}

	return NotifyGroup{Notifiers: notifiers, Timeout: cfg.NotifyTimeout}
}
