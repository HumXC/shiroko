package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"path"

	"github.com/HumXC/shiroko/logs"
	"github.com/HumXC/shiroko/server"
	"github.com/HumXC/shiroko/tools"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

var rootCommand *cobra.Command

func init() {
	rootCommand = &cobra.Command{
		Use: path.Base(os.Args[0]),
		RunE: func(cmd *cobra.Command, args []string) error {
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
			if err != nil {
				panic(err)
			}
			loglevel, err := flags.GetString("log-level")
			if err != nil {
				panic(err)
			}
			switch loglevel {
			case "debug":
				logs.SetLevel(slog.LevelDebug)
			case "info":
				logs.SetLevel(slog.LevelInfo)
			case "warn":
				logs.SetLevel(slog.LevelWarn)
			case "error":
				logs.SetLevel(slog.LevelError)
			default:
				return errors.New("log-level must be one of (debug|info|warn|error)")
			}
			mainRun(address, port, useDaemon)
			return nil
		},
	}
	flags := rootCommand.Flags()
	flags.StringP("address", "a", "0.0.0.0", "Listen address")
	flags.StringP("port", "p", "15600", "Listen port")
	flags.BoolP("daemon", "d", false, "Run as daemon")
	flags.StringP("log-level", "l", "info", "Log level (info|debug|warn|error)")
	rootCommand.AddCommand(&cobra.Command{
		Use:   "kill",
		Short: "kill all daemon",
		Run: func(cmd *cobra.Command, args []string) {
			Kill()
		},
	})
	rootCommand.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "list all process",
		Run: func(cmd *cobra.Command, args []string) {
			List()
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
		pid, err := Daemon()
		if err != nil {
			fmt.Println(err)
			return
		}
		if pid != 0 {
			_info, _ := json.MarshalIndent(map[string]any{
				"pid": pid,
			}, "", "    ")
			fmt.Println(string(_info))
			os.Exit(0)
		}
	}
	logs.WriteToStderr()
	lis, err := net.Listen("tcp", address+":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	serv := server.New()
	logs.Get("main").Info("grpc server will lieten to", "address", address, "port", port)
	serv.Serve(lis)
}
