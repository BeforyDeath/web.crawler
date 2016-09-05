package parser

import (
	"errors"
	"io"
	"io/ioutil"
	"net/url"

	"github.com/BeforyDeath/web.crawler/core"
	"github.com/BeforyDeath/web.crawler/storage"
	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/html"
)

func NormalizeItem(uri string) (i storage.Item, err error) {
	u, err := url.Parse(uri)
	if err != nil {
		return
	}

	if u.Scheme == "" {
		u.Scheme = core.Config.SchemeDefault
		if core.Config.Scheme != "" {
			u.Scheme = core.Config.Scheme
		}
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return i, errors.New("Scheme " + u.Scheme + ": not supported")
	}

	if u.Host == "" {
		if core.Config.Domain == "" {
			return i, errors.New("Host unknown")
		}
		u.Host = core.Config.Domain
	} else if core.Config.Domain != "" && u.Host != core.Config.Domain {
		return i, errors.New("Host outside: " + u.Host)
	} else {
		core.Config.Domain = u.Host
	}

	core.Config.Scheme = u.Scheme

	i.Url = u.String()

	return
}

func Reader(body io.ReadCloser) (b []byte, err error) {
	defer body.Close()
	b, err = ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	return
}

func Crawl(body io.ReadCloser) (items storage.ItemsType) {
	defer body.Close()

	z := html.NewTokenizer(body)

	items = make(storage.ItemsType)
	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			return
		case tt == html.StartTagToken:
			t := z.Token()

			isAnchor := false
			switch t.Data {
			case "a", "link", "script":
				isAnchor = true
			}

			if !isAnchor {
				continue
			}

			ok, link := getHref(t)
			if !ok {
				continue
			}

			item, err := NormalizeItem(link)
			if err != nil {
				log.Warn(err)
			} else {
				items.Add(item)
			}
		}
	}
	return
}

func getHref(t html.Token) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == "href" || a.Key == "src" {
			href = a.Val
			ok = true
		}
	}
	return
}
