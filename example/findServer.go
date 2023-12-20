package example

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/HumXC/shiroko/client"
)

func FindServer() ([]client.ShirokoServer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	ss, err := client.FindServer(ctx)
	if err != nil {
		return nil, err
	}
	if len(ss) == 0 {
		return nil, errors.New("No server found")
	}
	for _, s := range ss {
		fmt.Println("Name:", s.Name)
		fmt.Println("Model:", s.Model)
		fmt.Println("Addr:", s.Addr)
		fmt.Println()
	}
	return ss, nil
}
