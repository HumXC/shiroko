package main

func DefaultConfig() *Config {
	return &Config{
		Address:   "0.0.0.0",
		Port:      "15600",
		UseDaemon: false,
	}
}

type Config struct {
	Address   string
	Port      string
	UseDaemon bool
}
