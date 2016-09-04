package main

import (
	"net/http"
	"time"
	log "github.com/Sirupsen/logrus"
	"strings"
	"io"
)

var curl http.Client

func init() {
	curl = http.Client{
		Timeout: time.Duration(time.Second * 30),
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

	log.Info(res.Header.Get("Content-Type"))

	l.Status.Code = res.StatusCode
	l.Status.DataTime = time.Now()
	return res.Body, nil
}
