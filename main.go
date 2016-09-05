package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"

	"github.com/BeforyDeath/web.crawler/core"
	"github.com/BeforyDeath/web.crawler/parser"
	"github.com/BeforyDeath/web.crawler/storage"
	log "github.com/Sirupsen/logrus"
	"strings"
)

func main() {

	startUrl := flag.String("url", "", "Url link")
	flag.Parse()

	if *startUrl == "" {
		log.Fatal("Provide a url or domain: -url=arg")
	}

	item, err := parser.NormalizeItem(*startUrl)
	if err != nil {
		log.Fatal(err)
	}

	var path string = "download/" + core.Config.Domain

	err = os.MkdirAll(path+"/files/", 0777)
	if err != nil {
		log.Fatal(err)
	}

	storage.Items.Read(path + "/url.json")
	defer func() {
		log.Infof("Save %v urls", len(storage.Items))
		storage.Items.Save(path + "/url.json")
	}()

	storage.Items.Add(item)

	storage.Stack.Init()

	for {
		hash := storage.Stack.Pop()
		if hash == "" {
			log.Info("End")
			break
		}

		if item, ok := storage.Items[hash]; ok {
			body, err := storage.Curl(item)
			if err != nil {
				log.Error(err)
			}

			if body != nil {
				log.Infof("Edit :%v", item.Url)
				b, err := parser.Reader(body)
				if err != nil {
					log.Error(err)
					break
				}

				items := parser.Crawl(ioutil.NopCloser(bytes.NewBuffer(b)))
				storage.Items.Marge(items)

				fileGz := hash + storage.FileType[strings.Split(item.Status.ContentType, ";")[0]]
				err = storage.Gzip(b, path+"/files/"+fileGz+".gz")
				if err != nil {
					log.Error(err)
				}

			}
		}
	}
}
