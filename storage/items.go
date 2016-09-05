package storage

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
)

var Items = make(ItemsType)

type ItemsType map[string]*Item

type Item struct {
	Url    string
	Status status
}

type status struct {
	Code        int
	ContentType string
	DataTime    time.Time
}

var FileType = map[string]string{
	"text/html":              ".html",
	"application/javascript": ".js",
	"text/css":               ".css",
}

func (it ItemsType) Marge(items ItemsType) {
	for _, v := range items {
		if hash := it.Add(*v); hash != "" {
			Stack.Push(hash)
		}
	}
}

func (it ItemsType) Add(i Item) string {
	h := it.Hashed(i.Url)
	if _, ok := it[h]; !ok {
		it[h] = &i
		return h
	}
	return ""
}

func (it ItemsType) Hashed(str string) string {
	hashed := md5.New()
	hashed.Write([]byte(str))
	return hex.EncodeToString(hashed.Sum(nil))
}

func (it ItemsType) Read(filename string) error {
	log.Infof("Read JSON file: %v", filename)
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error(err)
		return err
	}
	if err = json.Unmarshal(f, &it); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (it ItemsType) Save(filename string) error {
	log.Infof("Save JSON file: %v", filename)
	fo, err := os.Create(filename)
	if err != nil {
		log.Error(err)
		return err
	}
	defer fo.Close()
	e := json.NewEncoder(fo)
	if err = e.Encode(it); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func Gzip(b1 []byte, path string) error {
	var b bytes.Buffer

	w := gzip.NewWriter(&b)
	w.Write(b1)
	w.Close()

	err := ioutil.WriteFile(path, b.Bytes(), 0666)
	if err != nil {
		return err
	}
	return nil
}
