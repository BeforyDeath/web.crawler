package main

import (
	"time"
	"crypto/md5"
	"encoding/hex"
	log "github.com/Sirupsen/logrus"
)

type Items map[string]*Link

type Link struct {
	Url    string
	Status status
}

type status struct {
	Code     int
	DataTime time.Time
}

func (s Items) Add(l string) string {
	h := s.Hashed(l)
	if _, ok := s[h]; !ok {
		s[h] = &Link{Url:l}
		log.Info(h)
		return h
	}
	return ""
}

func (s Items) Hashed(str string) string {
	hashed := md5.New()
	hashed.Write([]byte(str))
	return hex.EncodeToString(hashed.Sum(nil))
}