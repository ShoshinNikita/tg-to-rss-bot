package download

import (
	"net/url"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/rylio/ytdl"
)

const dataFolder = "data"

type Video struct {
	videInfo *ytdl.VideoInfo

	Title         string
	Description   string
	DatePublished time.Time
	Author        string
	LinkToAudio   string
}

func NewVideo(u *url.URL) (v Video, err error) {
	v = Video{}
	v.videInfo, err = ytdl.GetVideoInfo(u)
	if err != nil {
		return Video{}, err
	}
	v.Title = v.videInfo.Title
	v.Description = v.videInfo.Description
	v.DatePublished = v.videInfo.DatePublished
	v.Author = v.videInfo.Author
	v.LinkToAudio = dataFolder + "/" + v.videInfo.Title + ".mp3"

	return v, nil
}

func (v Video) Download() <-chan interface{} {
	responses := make(chan interface{}, 15)

	go func() {
		format, err := func() (form ytdl.Format, err error) {
			err = errors.New("there's no good format")
			for _, f := range v.videInfo.Formats {
				if f.Resolution == "" && f.Extension == "mp4" {
					// Can use any
					if form.AudioBitrate == 0 {
						form = f
						err = nil
						continue
					}
					// Choose the best AudioBitrate
					if f.AudioBitrate > form.AudioBitrate {
						form = f
					}
				}
			}
			return
		}()
		if err != nil {
			responses <- err
			close(responses)
			return
		}

		if _, err := os.Open(v.LinkToAudio); err == nil {
			responses <- errors.New("Video already exists")
			close(responses)
			return
		}

		f, err := os.Create(v.LinkToAudio)
		if err != nil {
			responses <- err
			close(responses)
			return
		}
		defer f.Close()

		responses <- "❕ Downloading..."

		err = v.videInfo.Download(format, f)
		if err != nil {
			// Delete bad file
			os.Remove(v.LinkToAudio)
			responses <- err
			close(responses)
			return
		}

		responses <- "✅ Downloading is finished"
		close(responses)
	}()

	return responses
}
