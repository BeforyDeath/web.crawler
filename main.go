package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"encoding/json"
	"os"
	"bytes"
	"io"
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

	Read("download/" + domain + "/list.json", &Storage)

	Storage.Add(u.String())

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
				BodyRead(boby, h)
			}

		} else {
			break
		}

		//Save("download/" + domain + "/list.json", Storage)
		//return
	}

	log.Info("end")
}

func BodyRead(boby io.ReadCloser, hash string) {
	defer boby.Close()

	b, err := ioutil.ReadAll(boby)
	if err == nil {

		boby_tmp := ioutil.NopCloser(bytes.NewBuffer(b))
		crawl(boby_tmp)
		Save("download/" + domain + "/list.json", Storage)

		err = SaveGzip(b, "download/" + domain + "/" + hash + ".gz")
		if err != nil {
			log.Error(err)
		}

	} else {
		log.Error(err)
	}
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

