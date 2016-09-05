package storage

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

var client http.Client

func init() {
	client = http.Client{
		Timeout: time.Duration(time.Second * 5),
	}
}

func Curl(i *Item) (io.ReadCloser, error) {
	res, err := client.Get(i.Url)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			i.Status.Code = http.StatusBadRequest
			return nil, err
		}
		if strings.Contains(err.Error(), "Client.Timeout") {
			i.Status.Code = http.StatusGatewayTimeout
			return nil, err
		}
		return nil, err
	}

	ContentType := res.Header.Get("Content-Type")
	i.Status.DataTime = time.Now()
	i.Status.ContentType = ContentType
	i.Status.Code = res.StatusCode

	if _, ok := FileType[strings.Split(ContentType, ";")[0]]; !ok {
		i.Status.Code = http.StatusNotExtended
		return nil, errors.New("Not supported: " + ContentType)
	}

	return res.Body, nil
}
