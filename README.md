# tg-to-rss-bot

Telegram bot, that adds received videos from YouTube into RSS feed.

**Requirements:**

- Docker

## Example of run script

Run `docker pull kirtis/tg-to-rss-bot:latest`

Use next script:

```bash
docker run --rm --name=tg-to-rss-bot \
-p 80:80 \
-v data:/app/data \
-v rss:/app/rss \
-e TOKEN=TELEGRAM_BOT_TOKEN \
-e HOST=https://SomeHost.com \
kirtis/tg-to-rss-bot:latest
```
