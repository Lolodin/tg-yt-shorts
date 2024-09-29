package converter

import (
	"context"
	"errors"
	"github.com/kkdai/youtube/v2"
	"io"
	"net/http"
	"time"
)

type Converter interface {
	Covert(url string) ([]byte, error)
}
type Video struct {
	client *youtube.Client
}

func NewVideo() *Video {
	return &Video{client: &youtube.Client{HTTPClient: http.DefaultClient}}
}

func (v *Video) Covert(url string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	vid, err := v.client.GetVideoContext(ctx, url)
	if err != nil {
		return nil, errors.Join(err, errors.New("cant get video from url: "+url))
	}
	format := getFormat(vid)
	stream, _, err := v.client.GetStreamContext(ctx, vid, format)
	defer stream.Close()
	if err != nil {
		return nil, err
	}
	buff, err := io.ReadAll(stream)
	if err != nil {
		return nil, err
	}

	return buff, nil
}

func getFormat(vid *youtube.Video) *youtube.Format {
	format := vid.Formats.WithAudioChannels()

	return &format[0]
}
