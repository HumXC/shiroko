package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/HumXC/shiroko/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	target := "192.168.3.204:15600"
	client, err := client.New(target, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatal(err)
	}
	info, err := client.Minicap.Info()
	if err != nil {
		log.Fatal(err)
	}
	err = client.Minicap.Stop()
	if err != nil {
		log.Fatal(err)
	}
	_ = client.Minicap.Start(info.Height, info.Width, info.Height, info.Width, info.Rotation)
	time.Sleep(1 * time.Second)
	reader, err := client.Minicap.Cat()
	if err != nil {
		log.Fatal(err)
	}
	buf := make([]byte, 1024*1024)
	for {
		time.Sleep(1 * time.Second)
		n, err := reader.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				continue
			} else {
				fmt.Println(err)
			}
		}
		fmt.Println(n)
	}

}
