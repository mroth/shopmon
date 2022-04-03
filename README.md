# 🏪 shopmon

Quick and dirty monitor bot to check for an in-stock product from any Shopify
powered store, and send a notification.


## Usage
Pre-built docker images are available on GitHub Container Registry for amd64,
arm64, and armhf architectures.

    docker run -d \
        -e DOMAIN="store.foo.com" \
        -e HANDLES="foobook-max-pro, foobook-mini" \
        -e RATE=60 \
        -e SLACK_WEBHOOK=https://hooks.slack.com/******** \
        -e DISCORD_WEBHOOK=https://discord.com/api/webhooks/******** \
        --name shopmon \
        ghcr.io/mroth/shopmon:main

See the config struct in `main.go` for the full list of supported configuration options.
