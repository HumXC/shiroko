package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/grandcat/zeroconf"
)

type ShirokoServer struct {
	Name, Model, Addr string
}

func FindServer(ctx context.Context) ([]ShirokoServer, error) {
	result := []ShirokoServer{}
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return nil, err
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			s := ShirokoServer{Name: entry.Instance}
			port := ""
			for _, t := range entry.Text {
				kv := strings.Split(t, "=")
				switch kv[0] {
				case "model":
					s.Model = kv[1]
				case "port":
					port = kv[1]
				}
			}
			s.Addr = fmt.Sprintf("%s:%s", entry.AddrIPv4[0], port)
			result = append(result, s)
		}
	}(entries)
	err = resolver.Browse(ctx, "_shiroko._tcp", "local.", entries)
	if err != nil {
		return nil, err
	}
	<-ctx.Done()
	// Wait some additional time to see debug messages on go routine shutdown.
	return result, nil
}
