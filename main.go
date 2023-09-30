package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/HumXC/shiroko/config"
	"github.com/HumXC/shiroko/server"
)

func main() {
	cfg := config.Get()
	if cfg.UseDaemon {
		daemon, err := Daemon()
		if err != nil {
			fmt.Println(err)
			return
		}
		if !daemon.IAmDaemon {
			info := map[string]any{
				"pid": daemon.Pid,
			}
			_info, _ := json.MarshalIndent(info, "", "    ")
			fmt.Println(string(_info))
			os.Exit(0)
		}
	}

	lis, err := net.Listen("tcp", cfg.Address+":"+cfg.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	serv := server.New()
	serv.Serve(lis)
}
