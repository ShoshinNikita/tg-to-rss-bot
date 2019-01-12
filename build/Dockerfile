FROM alpine
RUN apk update && \
	apk add ca-certificates ffmpeg && \
	update-ca-certificates
RUN mkdir /app
WORKDIR /app
COPY tg-to-rss-bot .
ENTRYPOINT [ "./tg-to-rss-bot" ]