package main

import "time"

type storage map[string]link

type link struct {
	Url    string
	status status
}

type status struct {
	Code     int
	DataTime time.Time
}