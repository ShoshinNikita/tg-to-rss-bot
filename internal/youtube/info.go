package youtube

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"unicode"

	"github.com/pkg/errors"
)

const (
	videoInfoURL = "https://www.youtube.com/get_video_info?&video_id=%s"
	thumbnailURL = "https://img.youtube.com/vi/%s/0.jpg"
)

func getVideoInfo(id string) (*Video, error) {
	stringQuery, err := fetchMeta(id)
	if err != nil {
		return nil, errors.Wrap(err, "can't fetch meta info")
	}

	query, err := url.ParseQuery(stringQuery)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse meta info")
	}

	v := &Video{}
	v.Author = query.Get("author")
	v.Title = query.Get("title")
	v.Filename = transformFilename(v.Title) + ".mp3"
	v.Description = v.Author + " - " + v.Title
	v.ThumbnailURL = fmt.Sprintf(thumbnailURL, id)

	formats := strings.Split(query.Get("url_encoded_fmt_stream_map"), ",")
	for _, f := range formats {
		values, err := url.ParseQuery(f)
		if err != nil {
			continue
		}

		// We need only url with itag == "18"
		itag := values.Get("itag")
		if itag != "18" {
			continue
		}

		videoURL := values.Get("url")
		videoURL += "&signature" + values.Get("sig")

		v.downloadURL = videoURL
	}

	if v.downloadURL == "" {
		return nil, errors.New("can't get link for video")
	}

	return v, nil
}

func fetchMeta(id string) (string, error) {
	url := fmt.Sprintf(videoInfoURL, id)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// transformFilename remove non-letter and non-digit runes and replace spaces with '-'
func transformFilename(filename string) string {
	res := make([]rune, 0, len(filename))

	for _, r := range filename {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			res = append(res, unicode.ToLower(r))
		case r == ' ':
			res = append(res, '-')
		}
	}

	return string(res)
}
