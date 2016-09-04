package main

import (
	"net/http"
	"time"
	log "github.com/Sirupsen/logrus"
	"strings"
	"io"
	"errors"
)

var curl http.Client

func init() {
	curl = http.Client{
		Timeout: time.Duration(time.Second * 5),
	}
	log.Info("init curl")
}

func Get(l *Link) (io.ReadCloser, error) {
	res, err := curl.Get(l.Url)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			l.Status.Code = http.StatusBadRequest
			return nil, err
		}
		if strings.Contains(err.Error(), "Client.Timeout") {
			l.Status.Code = http.StatusGatewayTimeout
			return nil, err
		}
		return nil, err
	}

	l.Status.ContentType = res.Header.Get("Content-Type")
	l.Status.Code = res.StatusCode
	l.Status.DataTime = time.Now()

	if !strings.Contains(l.Status.ContentType, "html") && !strings.Contains(l.Status.ContentType, "javascript") && !strings.Contains(l.Status.ContentType, "css") {
		l.Status.Code = http.StatusNotExtended
		return nil, errors.New("Not supported: " + l.Status.ContentType)
	}

	return res.Body, nil
}
