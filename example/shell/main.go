package main

import (
	"fmt"

	"github.com/HumXC/shiroko/client"
	"github.com/HumXC/shiroko/example"
)

func main() {
	ss, err := example.FindServer()
	if err != nil {
		fmt.Println(err)
		return
	}
	shiroko, err := client.New(ss[0].Addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	keys := []string{
		"ro.product.manufacturer",
		"ro.build.version.release",
		"ro.build.version.sdk",
		"ro.product.model",
	}
	for _, key := range keys {
		GetProp(shiroko, key)
	}
	_, err = shiroko.Shell.Run("input tap 10 10", 0)
	if err != nil {
		fmt.Println(err)
	}
}

func GetProp(shiroko *client.Client, key string) {
	p, err := shiroko.Shell.Getprop(key)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(key, p)
}
