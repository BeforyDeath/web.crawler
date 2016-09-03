package main

import (
	"net/http"
	"time"
	"strings"
	log "github.com/Sirupsen/logrus"
	"net/url"
	"io/ioutil"
	"compress/gzip"
	"os"
)

var curl http.Client

func _main() {

	curl = http.Client{
		Timeout: time.Duration(time.Second * 5),
	}

	u, err := url.Parse("/beforydeath.ru")
	//u, err := url.Parse("http://www.intel.com/content/dam/www/public/us/en/documents/manuals/64-ia-32-architectures-software-developer-manual-325462.pdf")
	log.Println(err)
	log.Println(u.Host)
	log.Println(u.ForceQuery)
	log.Println(u.Opaque)
	log.Println(u.Path)
	log.Println(u.RawPath)
	log.Println(u.RawQuery)
	log.Println(u.Scheme)
	log.Println(u.String())

	// todo ftp ...
	if u.Scheme != "http" && u.Scheme != "https" {

	}

	if u.Scheme == "" {
		u.Scheme = "http"
	}

	//u.Host = "beforydeath.ru"
	log.Println(u.String())
	return

	log.Println(get(u.String()))


}

func get(url string) int {
	res, err := curl.Get(url)
	defer res.Body.Close()

	log.Println(res.Header.Get("Content-Type"))

	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			return 0
		}
		if strings.Contains(err.Error(), "Client.Timeout") {
			log.Println(err)
			return 503
		}
		log.Println(err)
		return 1
	}



	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}

	writer, err := os.Create("output.txt")
	if err != nil {
		log.Println(err)
	}
	defer writer.Close()

	w := gzip.NewWriter(writer)
	w.Write(body)
	defer w.Close()

	return 1
}