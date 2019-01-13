# tg-to-rss-bot

Telegram bot, that adds received videos from YouTube into RSS feed.

**Requirements:**

- Docker

## Env vars

| Var   | Default | Required | Usage                                                      |
| ----- | ------- | -------- | ---------------------------------------------------------- |
| TOKEN | ""      | yes      | Telegram bot token                                         |
| HOST  | ""      | yes      | Host for creating links in RSS feed ($HOST/data/audio.mp3) |
| TLS   | true    | no       | If it is true, server will use `https` (default behavior)  |

## Example of run script

- Generate self-signed TLS certificate in `ssl` folder:

  `openssl req -x509 -nodes -newkey rsa:2048 -sha256 -keyout key.key -out cert.cert`

- Run `docker pull kirtis/tg-to-rss-bot:latest`

- Use next script:

  ```bash
  docker run --rm --name=tg-to-rss-bot \
  -p 80:80 \
  -v ./data:/app/data \
  -v ./rss:/app/rss \
  -v ./ssl:/app/ssl \
  -e TOKEN=TELEGRAM_BOT_TOKEN \
  -e HOST=https://SomeHost.com \
  -e TLS=true
  kirtis/tg-to-rss-bot:latest
  ```
