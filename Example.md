# Example of run script

```bat
docker run --rm --name=tg-to-rss-bot \
-p 600:80 \
-v data:/app/data \
-v rss/:/app/rss \
-e TOKEN=SOME_TOKEN \
-e HOST=https://SomeHost.com \
kirtis/tg-to-rss-bot
```