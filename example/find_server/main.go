package main

import (
	"context"
	"fmt"
	"time"

	"github.com/HumXC/shiroko/client"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	ss, err := client.FindServer(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(ss) == 0 {
		fmt.Println("No server found")
		return
	}
	for _, s := range ss {
		fmt.Println("Name:", s.Name)
		fmt.Println("Model:", s.Model)
		fmt.Println("Addr:", s.Addr)
		fmt.Println()
	}
}
