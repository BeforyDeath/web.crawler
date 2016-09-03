package main

import (
	"time"
)

type storage map[string]link

type link struct {
	Url      string
	Status   int
	DataTime time.Time
}
