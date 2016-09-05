package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/BeforyDeath/web.crawler/core"
	"github.com/BeforyDeath/web.crawler/parser"
	"github.com/BeforyDeath/web.crawler/storage"
	log "github.com/Sirupsen/logrus"
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
	var gzFile int
	defer func() {
		log.Infof("Save urls \t%v", len(storage.Items))
		log.Infof("Stack urls \t%v", len(storage.Stack.Nodes))
		log.Infof("Save gz \t%v", gzFile)
		storage.Items.Save(path + "/url.json")
	}()

	storage.Items.Add(item)

	storage.Stack.Init()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)

	go func() {
		for {
			hash := storage.Stack.Pop()
			if hash == "" {
				log.Info("Stack End")
				close(sigChan)
				return
			}

			if item, ok := storage.Items[hash]; ok {
				body, err := storage.Curl(item)
				if err != nil {
					log.Warn(err)
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
					} else {
						gzFile++
					}

				}
			}
		}
	}()

	for {
		select {
		case <-sigChan:
			log.Info("Exit web.crawler")
			return
		}
	}
}
