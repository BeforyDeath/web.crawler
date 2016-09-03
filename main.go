package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"encoding/json"
	"net/url"
	"errors"
)

func main() {

	link := flag.String("link", "", "Url link")
	flag.Parse()

	if *link == "" {
		log.Error("Provide a url or domain: -link=arg")
	}

	l, err := NormalizeLink(*link)
	if err != nil {
		log.Error(err)
	}
	log.Info(l)

	return

	//Stack := stack{}

	Storage := make(storage)
	Read("store.json", &Storage)

	log.Info("end")
}

func NormalizeLink(l string) (link string, err error) {
	u, err := url.Parse(l)

	if u.Scheme == "" {
		u.Scheme = "http"
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return link, errors.New("Scheme not http")
	}

	if u.Host == "" {
		return link, errors.New("Host nil")
	}
	link = u.String()
	return
}

func Read(filename string, st interface{}) error {
	log.Infof("Read JSON file: %v", filename)
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error(err)
		return err
	}
	if err = json.Unmarshal(f, &st); err != nil {
		log.Error(err)
		return err
	}
	return nil
}