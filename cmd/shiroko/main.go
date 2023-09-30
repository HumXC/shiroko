package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"path"

	"github.com/HumXC/shiroko/server"
	"github.com/HumXC/shiroko/tools"
	"github.com/spf13/cobra"
)

var rootCommand *cobra.Command

func init() {
	cfg := DefaultConfig()
	rootCommand = &cobra.Command{
		Use: path.Base(os.Args[0]),
		Run: func(cmd *cobra.Command, args []string) {
			mainRun(*cfg)
		},
	}
	flags := rootCommand.Flags()
	flags.StringVarP(&cfg.Address, "address", "a", cfg.Address, "Listen address")
	flags.StringVarP(&cfg.Port, "port", "p", cfg.Port, "Listen port")
	flags.BoolVarP(&cfg.UseDaemon, "daemon", "d", cfg.UseDaemon, "Run as daemon")
}
func main() {
	tools.Init(rootCommand)
	err := rootCommand.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func mainRun(cfg Config) {
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
