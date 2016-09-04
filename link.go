package main

import (
	"errors"
	"net/url"
	"golang.org/x/net/html"
	"io"
	log "github.com/Sirupsen/logrus"
	"os"
	"compress/gzip"
)

var domain string
var scheme string

func SchemeLink(link string) (u *url.URL, err error) {
	u, err = url.Parse(link)
	if err != nil {
		return
	}

	if u.Scheme == "" {
		u.Scheme = "http"
		if scheme != "" {
			u.Scheme = scheme
		}
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return u, errors.New("Scheme not http")
	}

	return
}

func HostLink(u *url.URL) error {
	if u.Host == "" {
		return errors.New("Host nil")
	}
	return nil
}

func NormalizeLink(link string) (u *url.URL, err error) {
	u, err = SchemeLink(link)
	if err != nil {
		return
	}

	err = HostLink(u)
	if err != nil {
		if domain != "" {
			u.Host = domain
			return u, nil
		}
		return
	}

	if u.Host != domain {
		return nil, errors.New("Gone beyond host")
	}

	return
}

func getHref(t html.Token) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
		if a.Key == "src" {
			href = a.Val
			ok = true
		}
	}
	return
}

func crawl(b io.ReadCloser) {
	defer b.Close()

	z := html.NewTokenizer(b)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			return
		case tt == html.StartTagToken:
			t := z.Token()

			isAnchor := false
			switch t.Data {
			case "a":
				isAnchor = true
			case "link":
				isAnchor = true
			case "script":
				isAnchor = true
			}

			if !isAnchor {
				continue
			}

			ok, l := getHref(t)
			if !ok {
				continue
			}

			u, err := NormalizeLink(l)
			if err != nil {
				//log.Error(err)
			} else {
				h := Storage.Add(u.String())
				if h != "" {
					Stack.Push(h)
					log.Info("add " + u.String())
				}

			}
		}
	}

}

func SaveGzip(b []byte, path string) error {

	writer, err := os.Create(path)
	if err != nil {
		return err
	}
	defer writer.Close()

	w := gzip.NewWriter(writer)
	w.Write(b)
	defer w.Close()

	return nil
}