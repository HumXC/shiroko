package config

import "flag"

var address = "localhost"
var port = "15600"
var useDaemon = false

func init() {
	flag.StringVar(&address, "a", address, "Listen Address")
	flag.StringVar(&port, "p", port, "Listen Port")
	flag.BoolVar(&useDaemon, "d", useDaemon, "Daemonize")
	flag.Parse()
}

type Config struct {
	Address   string
	Port      string
	UseDaemon bool
}

func Get() Config {
	return Config{
		Address:   address,
		Port:      port,
		UseDaemon: useDaemon,
	}
}
