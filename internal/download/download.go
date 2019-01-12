package download

import (
	"net/url"
	"time"
)

const dataFolder = "data"

type Video struct {
	Title         string
	Description   string
	DatePublished time.Time
	Author        string
	LinkToAudio   string
}

func NewVideo(u *url.URL) (v Video, err error) {
	return Video{}, nil
}

func (v Video) Download() <-chan interface{} {
	responses := make(chan interface{}, 15)
	close(responses)

	return responses
}
