package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/HumXC/shiroko/client"
)

const ADDRESS = "192.168.3.252:15600"

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
	client, err := client.New(ss[0].Addr)
	if err != nil {
		log.Fatal(err)
	}
	tools, err := client.Manager.List()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("所有工具: ", tools)
	for _, name := range tools {
		err := client.Manager.Health(name)
		if err != nil {
			fmt.Printf("工具[%s]不可用: %s\n", name, err)
		} else {
			fmt.Printf("工具[%s]可用\n", name)
		}
	}
	png, err := client.Screencap.Png("")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("截取设备屏幕到: screenshot.png")
	_ = os.WriteFile("screenshot.png", png, 0755)
}
