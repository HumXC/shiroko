package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"path"

	"github.com/HumXC/shiroko/logs"
	"github.com/HumXC/shiroko/server"
	"github.com/HumXC/shiroko/tools"
	"github.com/spf13/cobra"
)

var rootCommand *cobra.Command

func init() {
	rootCommand = &cobra.Command{
		Use: path.Base(os.Args[0]),
		Run: func(cmd *cobra.Command, args []string) {
			flags := cmd.Flags()
			address, err := flags.GetString("address")
			if err != nil {
				panic(err)
			}
			port, err := flags.GetString("port")
			if err != nil {
				panic(err)
			}
			useDaemon, err := flags.GetBool("daemon")
			mainRun(address, port, useDaemon)
			if err != nil {
				panic(err)
			}
		},
	}
	flags := rootCommand.Flags()
	flags.StringP("address", "a", "0.0.0.0", "Listen address")
	flags.StringP("port", "p", "15600", "Listen port")
	flags.BoolP("daemon", "d", false, "Run as daemon")
	rootCommand.AddCommand(&cobra.Command{
		Use:   "kill",
		Short: "kill all daemon",
		Run: func(cmd *cobra.Command, args []string) {
			Kill()
		},
	})

}
func main() {
	tools.Init(rootCommand)
	err := rootCommand.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func mainRun(address, port string, useDaemon bool) {
	if useDaemon {
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
	lis, err := net.Listen("tcp", address+":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	serv := server.New()
	logs.Get().Info("grpc server will lieten to", "address", address, "port", port)
	serv.Serve(lis)
}
