package core

var Config config

type config struct {
	Scheme        string
	SchemeDefault string
	Domain        string
}

func init() {
	Config.SchemeDefault = "http"
}
