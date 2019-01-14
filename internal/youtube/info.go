package youtube

import (
	"unicode"

	"github.com/knadh/go-get-youtube/youtube"
)

func getVideoInfo(id string) (*Video, error) {
	video, err := youtube.Get(id)
	if err != nil {
		return nil, err
	}

	return &Video{
		Author:      video.Author,
		Title:       video.Title,
		Filename:    transformFilename(video.Title) + ".mp3",
		Description: video.Author + " - " + video.Title,
		video:       &video,
	}, nil
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
