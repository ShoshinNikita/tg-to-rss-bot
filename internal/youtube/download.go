package youtube

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/ShoshinNikita/log"
	"github.com/ShoshinNikita/tg-to-rss-bot/internal/params"
	"github.com/knadh/go-get-youtube/youtube"
)

// init creates params.DataFolder
func init() {
	err := os.Mkdir(params.DataFolder, 0666)
	if err != nil && !os.IsExist(err) {
		log.Fatalf("can't create folder %s: %s", params.DataFolder, err)
	}
}

type Video struct {
	Author      string
	Title       string
	Filename    string
	Description string

	video *youtube.Video
}

func NewVideo(id string) (*Video, error) {
	return getVideoInfo(id)
}

// Message is used in Download function
type Message struct {
	Msg string

	IsFinished   bool
	IsFatalError bool
}

// Download downloads video.
// The last Message has IsFinished == true or IsFatalError == true
func (v *Video) Download() <-chan Message {
	url := v.video.Formats[0].Url
	results := make(chan Message)

	go func() {
		defer close(results)

		tempFilename := params.DataFolder + "temp-" + v.Filename
		videoFile, err := os.Create(tempFilename)
		if err != nil {
			results <- Message{
				Msg:          "can't create temp video file: " + err.Error(),
				IsFatalError: true,
			}
			return
		}

		// Get video content length
		contentLength, ok := func() (int64, bool) {
			resp, err := http.Head(url)
			if err != nil || resp.StatusCode == 403 || resp.Header.Get("Content-Length") == "" {
				return 0, false
			}

			header := resp.Header.Get("Content-Length")

			r, err := strconv.ParseInt(header, 10, 64)
			if err != nil {
				return 0, false
			}

			return r, true
		}()
		if !ok {
			results <- Message{
				Msg:          "can't define content length",
				IsFatalError: true,
			}
			return
		}

		resp, err := http.Get(url)
		if err != nil {
			results <- Message{
				Msg:          "request failed: " + err.Error(),
				IsFatalError: true,
			}
			return
		}
		defer resp.Body.Close()

		// Send status of downloading every 2 seconds
		go func() {
			ticker := time.NewTicker(time.Second * 2)
			defer ticker.Stop()

			var (
				offset   int64
				percents int64
				err      error
			)

			for range ticker.C {
				offset, err = videoFile.Seek(0, io.SeekCurrent)
				if err != nil {
					break
				}

				percents = 100 * offset / contentLength
				results <- Message{Msg: fmt.Sprintf("%d%%", percents)}

				if offset >= contentLength {
					break
				}
			}
		}()

		if _, err = io.Copy(videoFile, resp.Body); err != nil {
			results <- Message{
				Msg:          "can't download video file: " + err.Error(),
				IsFatalError: true,
			}
			return
		}

		videoFile.Close()

		// Delete file after converting
		defer func() {
			err := os.Remove(tempFilename)
			if err != nil {
				log.Errorf("can't delete video file %s: %s\n", tempFilename, err)
			}
		}()

		results <- Message{Msg: "Converting..."}

		ffmpeg, err := exec.LookPath("ffmpeg")
		if err != nil {
			results <- Message{
				Msg:          "can't find ffmpeg: " + err.Error(),
				IsFatalError: true,
			}
			return
		}

		cmd := exec.Command(ffmpeg, "-y", "-loglevel", "quiet", "-i", tempFilename, "-vn", params.DataFolder+v.Filename)
		if err := cmd.Run(); err != nil {
			results <- Message{
				Msg:          "ffmpeg exit with error: " + err.Error(),
				IsFatalError: true,
			}
			return
		}

		results <- Message{Msg: "Done", IsFinished: true}
	}()

	return results
}
