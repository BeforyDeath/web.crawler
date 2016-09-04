package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"encoding/json"
	"os"
	"bytes"
)

var Storage = make(Items)
var Stack stack

func main() {

	startLink := flag.String("link", "", "Url link")
	flag.Parse()

	if *startLink == "" {
		log.Fatal("Provide a url or domain: -link=arg")
	}

	u, err := SchemeLink(*startLink)
	if err != nil {
		log.Fatal(err)
	}

	err = HostLink(u)
	if err != nil {
		log.Fatal(err)
	}
	scheme = u.Scheme
	domain = u.Host

	err = os.MkdirAll("download/" + domain, 0777)
	if err != nil {
		log.Error(err)
	}

	//Read("download/" + domain + "/list.json", &Storage)

	Storage.Add(u.String())
	log.Info(Storage)

	for k, v := range Storage {
		if v.Status.Code == 0 {
			Stack.Push(k)
		}
	}
	log.Info(Stack)

	for {
		if h := Stack.Pop(); h != "" {
			log.Info(h)

			boby, err := Get(Storage[h])
			if err != nil {
				log.Error(err)
			}
			log.Info(Storage[h].Url)

			if boby != nil {

				b, err := ioutil.ReadAll(boby)
				if err == nil {
					boby1 := ioutil.NopCloser(bytes.NewBuffer(b))
					crawl(boby1)

					err = SaveGzip(b, "download/" + domain + "/" + h + ".gz")
					if err != nil {
						log.Error(err)
					}

				} else {
					log.Error(err)
				}
			}

		} else {
			break
		}
	}

	Save("download/" + domain + "/list.json", Storage)
	log.Info("end")
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

func Save(filename string, st interface{}) error {
	log.Infof("Save JSON file: %v", filename)
	fo, err := os.Create(filename)
	if err != nil {
		log.Error(err)
		return err
	}
	defer fo.Close()
	e := json.NewEncoder(fo)
	if err = e.Encode(st); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

